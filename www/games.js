$(document).ready(function() {
  weeks = document.getElementById("weekselector");

  createWeeksPaginationBar(weeks, loadGames);
});

var gamesCache = [];

// loadGames fetches and loads games into the table for the given week and year.
//
// If the data has already been retrieved from the source, it loads the games from an in-memory cache.
// Otherwise an AJAX call retrieves the data from the server.
function loadGames(year, week) {
  return function() {
    if (gamesCache[week] != null) {
      renderGamesTable(gamesCache[week]);
      return
    }

    $.getJSON("http://localhost:61389/games?year="+year+"&week="+week, function(games) {
      gamesCache[week] = games;
      renderGamesTable(games);
    });
  };
}

// setGameTable takes a list of JSON objects representing NFL Pick-Em' Games and populates
// games table with the information.
function renderGamesTable(games) {
  var table = document.getElementById("games").getElementsByTagName("tbody")[0];

  while(table.rows.length > 0) {
    table.deleteRow(-1);
  }

  for (let g of games) {
    var row = table.insertRow(table.rows.length);

    var cell = row.insertCell(0);
    cell.innerHTML = g.date;

    cell = row.insertCell(1);
    cell.innerHTML = g.away.city + " " + g.away.nickname;

    cell = row.insertCell(2);
    cell.innerHTML = g.awayScore;

    cell = row.insertCell(3);
    cell.innerHTML = g.home.city + " " + g.home.nickname;

    cell = row.insertCell(4);
    cell.innerHTML = g.homeScore;
  }
}
