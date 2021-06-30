const profile = document.getElementById("profile")
const profileDropdown = document.getElementById("profile-dropdown")
profile.addEventListener('click', () => {
  profileDropdown.classList.toggle("hidden")
})
const first = document.getElementById("first")
const frame = document.getElementById("frame")
const menus = document.querySelectorAll(".menu")
menus.forEach(target => {
  target.addEventListener('click', function () {
    if (!first.classList.contains('hidden')) {
      first.classList.add('hidden')
    }
    frame.src = this.dataset.url
    frame.classList.remove('hidden')
  })
})
