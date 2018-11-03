$('#copyButton').hover(function () {
    $('#copyButton').tooltip({title: 'Copy'})
                    .tooltip('show');
});

$('#copyButton').mouseleave(function () {
    $('#copyButton').attr('data-original-title', 'Copy');
});

$('#copyButton').click(function () {
    // create a hidden text area
    hiddenText = document.createElement('textarea');
    hiddenText.value = document.getElementById("result").value;
    hiddenText.style.position = 'absolute';
    hiddenText.style.left = '-1000px';
    document.body.appendChild(hiddenText);
    
    // select the hidden text area and copy
    hiddenText.select();
    document.execCommand('copy');

    // remove the hidden text area
    document.body.removeChild(hiddenText);

    // update the tooltip to inform the user
    $('#copyButton').attr('data-original-title', 'Copied!').tooltip('show');
});