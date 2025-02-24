$(document).ready(function() {
    $.getJSON("colors.json", function(data) {
        window.colors = data;
    });
});