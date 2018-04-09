$('.ui.dropdown').dropdown();

var timeout = 10000;
var oldStats = {};

function pollStats() {
    $.getJSON('/stats', function (newStats) {
        // Set all countups
        countup('dayLookups', 24, oldStats.dayLookups, newStats.dayLookups);
        countup('monthLookups', 720, oldStats.monthLookups, newStats.monthLookups);
        oldStats = newStats;

        // Perform this action every 10 seconds
        setTimeout(pollStats, timeout); // Poll stats every 10 seconds and re-apply to UI
    });
}

// countup animates the passed id with a new counted up value
function countup(id, hours, from, to, prefix, suffix) {
    // if from isn't set yet, calculate a rough initial value 
    if (from == undefined) {
        var rps = ((to/hours)/60)/60; // Average requests per second
        from = to-(rps*(timeout/1000)); // Subtract 10 seconds worth
    }

    // Configure countup options
    var options = { 
        useEasing: false, 
        useGrouping: true, 
        separator: ',', 
        decimal: '.'
    };

    // Apply a prefix if one is passed
    if (prefix != undefined && prefix != '') {
        options.prefix = prefix;
    }

    // Apply a suffix if one is passed
    if (suffix != undefined && suffix != '') {
        options.suffix = suffix;
    }

    // Trigger the countup animation
    var count = new CountUp(id, from, to, 0, 10, options);
    if (!count.error) {
        count.start();
    } else {
        console.error(count.error);
    }
}

$(document).ready(function () {
    pollStats();
    $('#stats-form').on('submit', function (e) {
        e.preventDefault();
        var area = document.getElementById('area').value;
        var tag = document.getElementById('tag').value;

        // Verify the parameters were passed
        if (area === '' || tag === '') {
            return;
        }

        // Open the json response page
        if (area == 'xbl' || area == 'psn') {
            window.open('/stats/' + area + '/' + tag);
        }
        window.open('/stats/pc/' + area + '/' + tag);
    });
});