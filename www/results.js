document.addEventListener("DOMContentLoaded", function() {
  years = document.getElementById("yearselector");
  weeks = document.getElementById("weekselector");

  createYearsWeeksPaginationBar(years, weeks, loadResults);
});

var resultsCache = [];

// loadResults fetches and loads pick results into the table for the given week and year.
//
// If the data has already been retrieved from the source, it loads the reuslts from an in-memory cache.
// Otherwise, an AJAX call retrieves the data from the server.
//
// Parameters:
//    year - NFL schedule year
//    week - NFL schedule week
function loadResults(year, week) {
  if (resultsCache[week] != null) {
    renderResultsTable(resultsCache[week]);
    return;
  }

  var request = new XMLHttpRequest();
  request.open("GET", "http://localhost:61389/results?year=" + year + "&week=" + week, true);

  request.onload = function() {
    if (this.status >= 200 && this.status < 400) {
      var results = JSON.parse(this.response);
      resultsCache[week] = results;
      renderResultsTable(results);
    } else {
      // TODO: Handle error gracefully
    }
  };

  request.send();
}


// renderResultsTable takes a list of JSON objects representing NFL Pick-Em' Pick Results and populates
// the results table with the information.
//
// Parameters:
//    results - The set of pick results to render
function renderResultsTable(results) {
  // Clear out what was here before
  table = document.getElementById("results");

  table.deleteTHead();
  table.deleteTFoot();
  tbody = table.getElementsByTagName("tbody")[0];
  while (tbody.rows.length > 0) {
    tbody.deleteRow(-1);
  }

  // Grab the list of users for these picks
  users = [];
  if (results.length != 0) {
    for (let p of results[0].picks) {
      users.push({"name": p.user.firstName, "points": 0});
    }
  }

  // Render the table
  for (let r of results) {
    var row = tbody.insertRow(tbody.rows.length)
      var cell = row.insertCell(row.cells.length);
    cell.innerHTML = r.game.away.nickname + "/" + r.game.home.nickname;

    for (let p of r.picks) {
      var cell = row.insertCell(row.cells.length);
      cell.innerHTML = p.selection.nickname;

      if (p.points != 1) {
        cell.innerHTML += (" (" + p.points + ")");
      }

      if ((r.game.awayScore > r.game.homeScore && p.selection.nickname == r.game.away.nickname) || (r.game.homeScore > r.game.awayScore && p.selection.nickname == r.game.home.nickname)) { 
        cell.className += "success";
        updateUserPoints(users, p.user.firstName, p.points);
      } else {
        cell.className += "danger";
      }
    }
  }

  // Create the header and footer now that we've determined the order of everything
  var header = table.createTHead();
  var hrow = header.insertRow(0);
  var footer = table.createTFoot();
  var frow = footer.insertRow(0);

  var hcell = hrow.insertCell(hrow.cells.length); // empty cell for the game

  var fcell = frow.insertCell(frow.cells.lenght);
  fcell.innerHTML = "<b>Total</b>";

  for (let u of users) {
    hcell = hrow.insertCell(hrow.cells.length);
    hcell.innerHTML = "<b>" + u.name + "</b>";

    fcell = frow.insertCell(frow.cells.length);
    fcell.innerHTML = "<b>" + u.points + "</b>";
  }
}

// updateUserPoints updates the tracked point total for the given user for display in the results table
//
// Parameters:
//  users - array of tracked users
//  user - the user in question we are looking for
//  points - how many points to increment their total by
function updateUserPoints(users, user, points) {
  for (let u of users) {
    if (u.name == user) {
      u.points += points;
      break;
    }
  }
}
