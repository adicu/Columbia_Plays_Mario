$(function() {
    $('.key').click(function(e) {
        e.preventDefault();
        var key = $(this).data('key');
        var username = $('.name').val();
        $.post('/press', JSON.stringify({
            'key': key,
            'username': username
        }));
    });
});
