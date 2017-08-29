var currentUser = null;

document.addEventListener("DOMContentLoaded", function() {
  currentUser = state();
  configureNavbar(currentUser != null);

  years = document.getElementById("yearselector");
  weeks = document.getElementById("weekselector");

  createYearsWeeksPaginationBar(years, weeks, loadStandings);
});

function loadStandings(year, week) {
  var request = new XMLHttpRequest();
  request.open("GET", "/api/totals?year=" + year + "&week=" + week + "&kind=cumulative", true);

  request.onload = function() {
    if (this.status >= 200 && this.status < 400) {
      var totals = JSON.parse(this.response);
      render(totals);
    } else {
      // TODO: Handle error gracefully
    }
  };

  // TODO: Handle error request.onerror

  request.send();
}

// renderStandingsTable renders the standings table based on an array of totals in JSON form
// 
// Parameters:
//  weekTotals - JSON response from the pickem' API
function render(weekTotals) {
  // Clear out what was here before
  table = document.getElementById("standings").getElementsByTagName("tbody")[0];
  while (table.rows.length > 0) {
    table.deleteRow(-1);
  }

  // We're going to create an array of objects that look like
  // {
  //    "name": "",
  //    "totals": []
  // }
  //
  // Then we're going to sort them by the sum(points). That will give us enough to 
  // make our table.
  var standings = [];

  for (let t of weekTotals) {
    if (!containsUser(standings, t.user.firstName)) {
      standings.push({"name": t.user.firstName, "totals": []});
    }

    addTotalToStandings(standings, t.user.firstName, t.week, t.total);
  }

  standings.sort(function(a, b) {
    var aSum = 0;
    var bSum = 0;

    for (let t of a.totals) {
      aSum += t;
    }

    for (let t of b.totals) {
      bSum += t;
    }

    return bSum - aSum;
  });

  console.log(standings);

  var table = document.getElementById("standings").getElementsByTagName("thead")[0];

  while(table.rows.length > 0) {
    table.deleteRow(-1);
  }
  
  var row = table.insertRow(table.rows.length);
  var cell = row.insertCell(row.cells.length);
  cell.appendChild(document.createTextNode("Name"));

  for (var i=0; i < standings[0].totals.length; i++) {
    var cell = row.insertCell(row.cells.length);
    cell.appendChild(document.createTextNode(i+1));
  }

    cell = row.insertCell(row.cells.length);

  cell = row.insertCell(row.cells.length);
  cell.appendChild(document.createTextNode("Raw Total"));

  cell = row.insertCell(row.cells.length);
  cell.appendChild(document.createTextNode("Adjusted Total"));

  // Move to the tbody
  table = document.getElementById("standings").getElementsByTagName("tbody")[0];

  for (s of standings) {
    var row = table.insertRow(table.rows.length);

    var cell = row.insertCell(row.cells.length);
    cell.appendChild(document.createTextNode(s.name));

    for (var i=0; i < s.totals.length; i++) {
      var cell = row.insertCell(row.cells.length);
      cell.appendChild(document.createTextNode(s.totals[i]));
    }

    cell = row.insertCell(row.cells.length);

    cell = row.insertCell(row.cells.length);
    cell.appendChild(document.createTextNode(s.totals));

    cell = row.insertCell(row.cells.length);
    if (s.totals.length > 1) {
      cell.appendChild(document.createTextNode(adjustedTotalPoints(s.totals)));
    } else {
      cell.appendChild(document.createTextNode(totalPoints(s.totals)));
    }
  }
}

// totalPoints returns the sum of all week totals
//
// Parameters:
//  totals - array of week totals
function totalPoints(totals) {
  var total = 0;

  for (t of totals) {
    total += t;
  }

  console.log(total);

  return total;
}

// adjustedTotalPoints returns the sum minus the lowest week total
//
// Parameters:
//  totals - array of week totals
function adjustedTotalPoints(totals) {
  var min = 100;

  for (t of totals) {
    if (t < min) {
      min = t;
    }
    total += t;
  }

  console.log(total);

  return total - min;
}

// containsUser determines if the given user is contained in the array
//
// Parameters:
//  standings - array of objects
//  user - username to search for
function containsUser(standings, user) {
  for (s of standings) {
    if (s.name == user) {
      return true;
    }
  }

  return false;
}

// addTotalToStandings finds the user specified and appends the given total to list of totals
//
// Parameters:
//  standings - array of objects
//  user - user to add the total too
//  week - week the total represents
//  total - point total for the week
function addTotalToStandings(standings, user, week, total) {
  for (s of standings) {
    if (s.name == user) {
      s.totals[week-1] = total;
      return;
    }
  }
}
