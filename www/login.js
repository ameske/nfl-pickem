var currentUser = null;

document.addEventListener("DOMContentLoaded", function() {
  currentUser = state();

  configureNavbar(currentUser != null);
});

