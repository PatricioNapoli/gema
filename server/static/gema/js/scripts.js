$('body').mousemove(function(e){
    var amountMovedX = (e.pageX * -1 / 12);
    var amountMovedY = (e.pageY * -1 / 12);
    $(this).css('background-position', amountMovedX + 'px ' + amountMovedY + 'px');
});