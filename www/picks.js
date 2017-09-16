// Keep track of the current user
var currentUser = null;

document.addEventListener("DOMContentLoaded", function() {
  currentUser = state();

  configureNavbar(currentUser != null);

  years = document.getElementById("yearselector");
  weeks = document.getElementById("weekselector");

  submitButton = document.getElementById("submitpicks");
  submitButton.onclick = submitPicks;

  createYearsWeeksPaginationBar(years, weeks, loadPicks);
});

// Keep track of the original JSON that we got back from the server for the
// current week we are modifying. We'll need it to submit the picks back in
// a way the server can identify which pick to update.
var currentPicks = null;

// submitPicks extracts the current pick information from the table, updating
// our in memory view of the pick set. It then validates it, and submits it 
// to the server for storage if it passes validation.
function submitPicks() {
  table = document.getElementById("picks").getElementsByTagName("tbody")[0];

  // The picks are going to be in the same order in the table as they are in
  // the cached version of the JSON. We'll update "selection" and "points" for each.
  for (i=0; i < table.rows.length; i++) {
    let selection = table.rows[i].cells[3].firstChild.value;

    // Ignore picks that haven't been made that have locked (they will be empty)
    if (selection == "" || selection == null) {
      continue;
    }

    let points = table.rows[i].cells[4].firstChild.value;
    let split = selection.split(' ');

    if (split.length == 2) {
      currentPicks[i].selection.city = split[0];
      currentPicks[i].selection.nickname = split[1];
    } else {
      currentPicks[i].selection.city = split[0] + " " + split[1];
      currentPicks[i].selection.nickname = split[2];
    }

    currentPicks[i].points = parseInt(points);
  }

  if (!isValid(currentPicks)) {
    alert("Unable to send picks. Please correct your point values.");
    return;
  }

  let years = document.getElementById("yearselector");
  let year = currentlySelectedElementValue(years);
  let weeks = document.getElementById("weekselector");
  let week = currentlySelectedElementValue(weeks);

  let request = new XMLHttpRequest();
  request.open("POST", "/api/picks?year="+year+"&week="+week+"&username=" + currentUser.Username, true);
  request.withCredentials = true;
  request.setRequestHeader("Content-Type", "application/json");

  request.onload = function() {
    alert("Status: " + this.status + "\nResponse: " + this.response);
  };

  request.send(JSON.stringify(currentPicks));
} 

// isValid determines if a pick set uses a valid amount of special points
//
// Parameters:
//  picks - array of picks
function isValid(picks) {
  let seven = 0;
  let five = 0;
  let three = 0;

  for (p of picks) {
    switch (p.points) {
      case 7: seven++;
              break;
      case 5: five++;
              break;
      case 3: three++;
              break;
    }
  }

  if (seven > 1) {
    alert("Too many seven picks. You may only have 1, but you have " + seven);
  }

  if (five > 2) {
    alert("Too many five picks. You may only have 2, but you have " + five);
  }

  if (three > 5) {
    alert("Too many three picks. You may only have 5, but oyu have " + three);
  }

  return seven <= 1 && five <= 2 && three <= 5;
}

// loadPicks fetches and loads picks into the table for the given week and year.
//
// Parameters:
//    year - NFL schedule year
//    week - NFL schedule week
function loadPicks(year, week) {
  var request = new XMLHttpRequest();
  request.open("GET", "/api/picks?year="+year+"&week="+week+"&username=" + currentUser.Username, true);
  request.withCredentials = true;

  request.onload = function() {
    if (this.status >= 200 && this.status < 400) {
      var picks = JSON.parse(this.response);
      currentPicks = picks;
      render(picks);

      // unhide the submit button if it was hidden
      button = document.getElementById("submitpicks").removeAttribute("style");
    } else {
      // TODO: Handle error gracefully
      button = document.getElementById("submitpicks").removeAttribute("style");
      button.style = "display: none";
    }
  };

  // TODO: Handle error request.onerror

  request.send();
}

// render creates a picks table and loads it into the DOM
//
// Parameters:
//    root - The root of the DOM table element to render the picks list in
function render(picks) {
  var table = document.getElementById("picks").getElementsByTagName("tbody")[0];

  while(table.rows.length > 0) {
    table.deleteRow(-1);
  }

  var now = new Date();

  for (p of picks) {
    let row = table.insertRow(table.rows.length);
    renderPick(now, p, row)
  }
}

// renderPick renders a pick in the picks table at the given row
//
// Parameters:
//  now - the current time
//  pick - the pick to render
//  row - (HTMLRowElement) the row to render the pick in
function renderPick(now, pick, row) {
  let gametime = Date.parse(pick.game.date);

  let cell = row.insertCell(row.cells.length);
  cell.appendChild(document.createTextNode(pick.game.date));

  cell = row.insertCell(row.cells.length);
  cell.appendChild(document.createTextNode(pick.game.home.city + " " + pick.game.home.nickname));

  cell = row.insertCell(row.cells.length);
  cell.appendChild(document.createTextNode(pick.game.away.city + " " + pick.game.away.nickname));

  if (gametime < now) {
    cell = row.insertCell(row.cells.ength);
    cell.appendChild(document.createTextNode(pick.selection.city + " " + pick.selection.nickname));

    cell = row.insertCell(row.cells.length);
    cell.appendChild(document.createTextNode(pick.points));

  } else {
    cell = row.insertCell(row.cells.ength);
    cell.appendChild(renderTeamSelection(pick));

    cell = row.insertCell(row.cells.length);
    cell.appendChild(renderPointSelection(pick));
  }
}

// renderTeamSelection renders the HTML for a team selection element
//
// Parameter:
//  pick - the pick to render the team selection element for
function renderTeamSelection(pick) {
  let select = document.createElement("select");

  let option = document.createElement("option");
  option.value = pick.game.home.city + " " + pick.game.home.nickname;
  option.text = pick.game.home.city + " " + pick.game.home.nickname;
  option.selected = (pick.game.home.nickname == pick.selection.nickname);
  select.appendChild(option)

  option = document.createElement("option");
  option.value = pick.game.away.city + " " + pick.game.away.nickname;
  option.text = pick.game.away.city + " " + pick.game.away.nickname;
  option.selected = (pick.game.away.nickname == pick.selection.nickname);
  select.appendChild(option)

  return select;
}

// renderPointSelection renders the HTML for a point selection element
//
// Parameter:
//  pick - the pick to render the point selection element for
function renderPointSelection(pick) {
  let select = document.createElement("select");

  let option = document.createElement("option");
  option.value = 1;
  option.text = "1";
  option.selected = (pick.points == 1);
  select.appendChild(option);

  option = document.createElement("option");
  option.value = 3;
  option.text = "3";
  option.selected = (pick.points == 3);
  select.appendChild(option);

  option = document.createElement("option");
  option.value = 5;
  option.text = "5";
  option.selected = (pick.points == 5);
  select.appendChild(option);


  option = document.createElement("option");
  option.value = 7;
  option.text = "7";
  option.selected = (pick.points == 7);
  select.appendChild(option);

  return select;
}
