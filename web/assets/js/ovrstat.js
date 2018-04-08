$('.ui.dropdown').dropdown();

var oldStats = {};

function pollStats() {
    $.getJSON('/stats', function (newStats) {
        // Set all countups
        countup('dayLookups', oldStats.dayLookups, newStats.dayLookups);
        countup('monthLookups', oldStats.monthLookups, newStats.monthLookups);
        oldStats = newStats;

        // Perform this action every 10 seconds
        setTimeout(pollStats, 10000); // Poll stats every 10 seconds and re-apply to UI
    });
}

// countup animates the passed id with a new counted up value
function countup(id, from, to, prefix, suffix) {
    // if from isn't set yet, set it to the initial to value
    if (from == undefined) {
        from = to
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