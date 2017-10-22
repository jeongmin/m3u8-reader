document.addEventListener('DOMContentLoaded', function() {
    $(function () {
	$('[data-toggle="tooltip"]').tooltip()
    })

    $('#read').click(function() {
	m3u8Url = $('#m3u8Url').val().trim();
	readM3u8(m3u8Url, '#master-playlist');
    });

    $('#btn-clear').click(function() {
	$('#master-playlist').empty();
    });


    $('#master-playlist').on('click', '.variant-url', function() {
    	m3u8Url = $(this).attr('href').trim();
    	appendMediaPlaylist(m3u8Url, $(this).parent());
    });    
}, false);

function readM3u8(m3u8Url, selector) {
    $.ajax({
    	url: 'http://localhost:9000',
    	type: 'POST',
    	data: {
    	    m3u8Url: m3u8Url
    	},
    	dataType: 'text',
    	success: function (result) {    	    
	    $(selector).append(result);
    	},
    	error: function (xhr, status, errorThrown) {
    	    alert('failed');
    	}
    });
}

function appendMediaPlaylist(m3u8Url, element) {
    $.ajax({
    	url: 'http://localhost:9000',
    	type: 'POST',
    	data: {
    	    m3u8Url: m3u8Url
    	},
    	dataType: 'text',
    	success: function (result) {    	    
	    element.append(result);
    	},
    	error: function (xhr, status, errorThrown) {
    	    alert('failed');
    	}
    });
}
