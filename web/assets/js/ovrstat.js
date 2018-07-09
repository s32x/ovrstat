$('.ui.dropdown').dropdown();

$(document).ready(function () {
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