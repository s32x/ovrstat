$('.ui.dropdown').dropdown();

$(document).ready(function () {
    $('#stats-form').on('submit', function (e) {
        e.preventDefault();
        var platform = document.getElementById('platform').value;
        var tag = document.getElementById('tag').value;
        window.open('/stats/' + platform + '/' + tag);
    });
});