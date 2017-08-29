var currentUser = null;

$(document).ready(function() {
  currentUser = state();
  configureNavbar(currentUser != null);
  loadCurrentStandings();
});

function loadCurrentStandings() {
  $.getJSON("/api/totals?type=cumulative&year=" + currentYear + "&week=" + currentWeek, function(totals) {
          console.log(totals);
          standings = makeStandings(totals);

          var rows = [];
          var firstPlace = 0;

          if (standings.length != 0) {
            firstPlace = standings[0].total;
          }

          for (const row of standings) {
            rows.push( "<tr>" +
                       "<td>" + row.name + "</td>" +
                       "<td>" + row.total + "</td>" +
                       "<td>" + (firstPlace - row.total) + "</td>" + 
                       "</tr>");
          }

          $("#standings tbody").append(rows);
  });
}

/*
 * Process the raw totals and return an array of JSON elements representing the standings.
 *
 * Each row looks like:
 * {
 *    "name": <Name>,
 *    "total": <Total>,
 *    "pointsBack": <Points Back>
 * }
 *
 * These rows are sorted already for easy display in a standings table.
 */
function makeStandings(weekTotals) {
      var userTotals = {}

      for (const wt of weekTotals) {
            if (wt.user.firstName in userTotals) {
                userTotals[wt.user.firstName].total += wt.total;
                if (wt.total < userTotals[wt.user.firstName].min) {
                    userTotals[wt.user.firstName].min = wt.total;
                }
            } else {
                userTotals[wt.user.firstName] = {name: wt.user.firstName, total: wt.total, min: wt.total}
            }
      }

      var standings = [];
      for (ut in userTotals) {
          standings.push({name: userTotals[ut].name, total: userTotals[ut].total - userTotals[ut].min});
      }

      standings.sort(function(a, b) {
        if (a.total < b.total) {
          return 1;
        } else if (a.total > b.total) {
          return -1;
        } else {
          return 0;
        }
      });
      
      return standings
}
