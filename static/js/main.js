$(function() {
    $('.key').click(function(e) {
        e.preventDefault();
        var key = $(this).data('key');
        $.post('/press', JSON.stringify({'key': key}));
    });
});
