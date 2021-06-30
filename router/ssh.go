package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/n-creativesystem/go-fwncs"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/websocket"
)

const END_OF_TRANSMISSION = "\u0004"

type terminalMessage struct {
	Op, Data, SessionID string
	Rows, Cols          int
}

type consoleReader struct {
	dst     io.ReadCloser
	log     fwncs.ILogger
	session *ssh.Session
}

func (w *consoleReader) Read(p []byte) (int, error) {
	n, err := w.dst.Read(p)
	if err != nil {
		w.log.Error(err)
		return copy(p, END_OF_TRANSMISSION), err
	}
	if n > 0 {
		var msg terminalMessage
		if err := json.Unmarshal(p[:n], &msg); err != nil {
			w.log.Error(err)
			return copy(p, END_OF_TRANSMISSION), err
		}
		w.log.Info(fmt.Sprintf("%v", msg))
		switch msg.Op {
		case "stdin":
			w.log.Info(fmt.Sprintf("stdin: %s", msg.Data))
			return copy(p, []byte(msg.Data)), nil
		case "resize":
			w.log.Info(fmt.Sprintf("resize h: %d w: %d", msg.Rows, msg.Cols))
			if err := w.session.WindowChange(msg.Rows, msg.Cols); err != nil {
				w.log.Error(err)
				return copy(p, []byte(END_OF_TRANSMISSION)), err
			}
			return 0, nil
		default:
			return copy(p, []byte(END_OF_TRANSMISSION)), fmt.Errorf("unknown message tyoe '%s'", msg.Op)
		}
	}
	return 0, nil
}

func (w *consoleReader) Close() error {
	return w.dst.Close()
}

var _ io.ReadCloser = (*consoleReader)(nil)

func wrapReadCloser(dst io.ReadCloser, log fwncs.ILogger, session *ssh.Session) io.ReadCloser {
	return &consoleReader{dst: dst, log: log, session: session}
}

type consoleWriter struct {
	dst io.WriteCloser
	log fwncs.ILogger
}

func (w *consoleWriter) Write(p []byte) (n int, err error) {
	n, err = w.dst.Write(p)
	if n > 0 {
		w.log.Info(string(p[:n]))
	}
	return
}

func (w *consoleWriter) Close() error {
	return w.dst.Close()
}

var _ io.WriteCloser = (*consoleWriter)(nil)

func wrapWriteCloser(dst io.WriteCloser, log fwncs.ILogger) io.WriteCloser {
	return &consoleWriter{dst: dst, log: log}
}

func handlerSSH(c fwncs.Context) {
	host := c.QueryParam("host")
	port := c.QueryParam("port")
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		sshConfig := &ssh.ClientConfig{
			User: "kube",
			Auth: []ssh.AuthMethod{
				ssh.Password("kube"),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		client, err := ssh.Dial("tcp", net.JoinHostPort(host, port), sshConfig)
		if err != nil {
			c.Logger().Error(err)
			return
		}
		defer client.Close()
		session, err := client.NewSession()
		if err != nil {
			c.Logger().Error(err)
			return
		}
		defer session.Close()
		termmodes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		if err := session.RequestPty("xterm-256color", 80, 120, termmodes); err != nil {
			c.Logger().Error(err)
			return
		}
		combinedOut := wrapWriteCloser(ws, c.Logger())
		session.Stdout = combinedOut
		session.Stderr = combinedOut
		session.Stdin = wrapReadCloser(ws, c.Logger(), session)
		if err := session.Shell(); nil != err {
			c.Logger().Error(err)
			return
		}
		if err := session.Wait(); nil != err {
			c.Logger().Error(err)
		}
	}).ServeHTTP(c.Writer(), c.Request())
}
