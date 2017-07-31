// createWeeksPaginationBar creates a paging table to allow the user to select between the
// weeks of a given year for the page.
//
// Parameters:
//      root - The <UL> DOM object of the pager to be updated
//      onClick - a function of the form f(year, week) -> f() that will be called when a week is clicked
function createWeeksPaginationBar(root, onClick) {
  $.getJSON("http://localhost:61389/current", function(current) {
    while(root.hasChildNodes()) {
      root.removeChild(root.lastChild);
    }

    // Ignore this for now, we need to gracefully display all the seasons
    // somehow

    for (i=1; i <= 17; i++) {
      var a = document.createElement("A");
      a.setAttribute("href", "#");
      a.innerHTML = i;
      a.onclick = onClick(2017, i);

      var l = document.createElement("LI");
      l.appendChild(a);

      root.appendChild(l);
    }
  });
}

