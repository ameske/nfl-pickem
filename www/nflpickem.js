// yearOnClick is run every time an element of a yearselector nav bar
// is clicked.
//
// Parameters:
//      yearRoot - reference to the yearselector navbar
//      weekRoot - reference to the weekselector navbar
//      element - the element that was clicked
//      updateFn - the function to run that updates the rest of the page
function yearOnClickFunc(yearRoot, weekRoot, element, updateFn) {
  return function() {
    unselectAll(yearRoot);
    element.classList.add("active");

    week = currentlySelectedElementValue(weekRoot);

    updateFn(parseInt(element.innerText), week);
  }
}

// weekOnClickFunc is run every time an element of a weekselector nav bar
// is clicked.
//
// Parameters:
//      weekRoot - reference to the weekselector navbar
//      yearRoot - reference to the yearselector navbar
//      element - the element that was clicked
//      updateFn - the function to run that updates the rest of the page
function weekOnClickFunc(weekRoot, yearRoot, element, updateFn) {
  return function() {
    unselectAll(weekRoot);
    element.classList.add("active");

    year = currentlySelectedElementValue(yearRoot);

    updateFn(year, parseInt(element.innerText));
  }
}

// currentlySelectedElementValue returns the value of the element in the collection
// that has the active CSS class set. If there are multiple elements, the last value
// is used.
//
// Parameters:
//    root  - The collection
function currentlySelectedElementValue(root) {
  var value = -1;
  root.childNodes.forEach(function(val, idx, obj) {
    if (val.classList.contains("active")) {
      value = parseInt(val.innerText);
    }
  });

  return value;
}

// unselectAll removes the "active" class from all child elements.
//
// Parameters:
//      root - The root DOM object, assumed to be something with children
//              that appear "selected" if they contain the "active" class
function unselectAll(root) {
  root.childNodes.forEach(function(val, idx, obj) {
    val.classList.remove("active");
  });
}

// createWeeksPaginationBar creates a paging table to allow the user to select between the
// weeks of a given year for the page.
//
// Parameters:
//      root - The <UL> DOM object of the pager to be updated
//      onClick - a function of the form f(year, week) -> f() that will be called when a week is clicked
function createYearsWeeksPaginationBar(yearRoot, weekRoot, onClick) {
  var request = new XMLHttpRequest();
  request.open("GET", "/api/years", true);

  request.onload = function() {
    if (this.status >= 200 && this.status < 400) {
      var years = JSON.parse(this.response);
      while(yearRoot.hasChildNodes()) {
        yearRoot.removeChild(yearRoot.lastChild);
      }

      while(weekRoot.hasChildNodes()) {
        weekRoot.removeChild(weekRoot.lastChild);
      }

      for (i=0; i < years.years.length; i++) {
        var a = document.createElement("A");
        a.setAttribute("href", "#");
        a.innerHTML = years.years[i];

        var l = document.createElement("LI");
        if (i == years.years.length - 1) {
          l.classList.add("active");
          currentlySelectedYear = years.years[i];
        }
        l.onclick = yearOnClickFunc(yearRoot, weekRoot, l, onClick);

        l.appendChild(a);

        yearRoot.appendChild(l);
      }

      for (i=1; i <= 17; i++) {
        var a = document.createElement("A");
        a.setAttribute("href", "#");
        a.innerHTML = i;

        var l = document.createElement("LI");
        l.onclick = weekOnClickFunc(weekRoot, yearRoot, l, onClick);

        l.appendChild(a);

        weekRoot.appendChild(l);
      }
    } else {
      // TODO: Handle error gracefully
    }
  };

  // TODO: Handle error request.onerror

  request.send();

}
