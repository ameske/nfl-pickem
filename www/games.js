var currentUser = null;

document.addEventListener("DOMContentLoaded", function() {
  currentUser = state();
  configureNavbar(currentUser != null);

  years = document.getElementById("yearselector");
  weeks = document.getElementById("weekselector");

  createYearsWeeksPaginationBar(years, weeks, loadGames);
});

var gamesCache = [];

// loadGames fetches and loads games into the table for the given week and year.
//
// If the data has already been retrieved from the source, it loads the games from an in-memory cache.
// Otherwise an AJAX call retrieves the data from the server.
//
// Parameters:
//    year - NFL schedule year
//    week - NFL schedule week
function loadGames(year, week) {
  if (gamesCache[week] != null) {
    renderGamesTable(gamesCache[week]);
    return
  }

  var request = new XMLHttpRequest();
  request.open("GET", "/api/games?year="+year+"&week="+week, true);

  request.onload = function() {
    if (this.status >= 200 && this.status < 400) {
      var games = JSON.parse(this.response);
      gamesCache[week] = games;
      renderGamesTable(games);
    } else {
      // TODO: Handle error gracefully
    }
  };

  // TODO: Handle error request.onerror

  request.send();
}

// setGameTable takes a list of JSON objects representing NFL Pick-Em' Games and populates
// games table with the information.
//
// Parameters:
//    root - The root of the DOM table element to render the games list in
function renderGamesTable(games) {
  var table = document.getElementById("games").getElementsByTagName("tbody")[0];

  while(table.rows.length > 0) {
    table.deleteRow(-1);
  }

  for (let g of games) {
    var row = table.insertRow(table.rows.length);

    var cell = row.insertCell(row.cells.length);
    cell.innerHTML = g.date;

    cell = row.insertCell(row.cells.length);
    cell.innerHTML = g.away.city + " " + g.away.nickname;

    cell = row.insertCell(row.cells.length);
    cell.innerHTML = g.awayScore;

    cell = row.insertCell(row.cells.length);
    cell.innerHTML = g.home.city + " " + g.home.nickname;

    cell = row.insertCell(row.cells.length);
    cell.innerHTML = g.homeScore;
  }
}
