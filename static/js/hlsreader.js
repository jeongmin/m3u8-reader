var serverUrl = 'http://220.230.118.50:8080';
//var serverUrl = 'http://localhost:8080';
var apiRawData = '/rawdata';

$(document).ready(function() {
    $("input[name=m3u8Url]").keydown(function (key) {
	// handle Enter key to read a url in an input tag.
	if (key.keyCode == 13) {
	    read();
	}
    });
    
    $(function () {
	$('[data-toggle="tooltip"]').tooltip();
    })

    $('#read').click(function() {
	read();
    });

    function read() {
	m3u8Url = $('#m3u8Url').val().trim();
	readM3u8(m3u8Url, $('#master-playlist'));
    }

    $('#master-playlist').on('click', '.rawdata', function() {
	if ($(this).children('span').attr('class').match('need-moredata') == null) {
	    return;
	}
	
	apiUrl      = serverUrl + apiRawData;
	m3u8Url     = $(this).attr('data-raw-url').trim();
	targetId    = $(this).attr('data-target-id');
	badge       = $(this).children('span');
	placeHolder = $(this).parent();
	originalBadgeText = badge.text();
	badge.text('downloading');
	$.ajax({
    	    url: apiUrl,
    	    type: 'POST',
    	    data: {
    		m3u8Url: m3u8Url
    	    },
    	    dataType: 'text',
    	    success: function (result) {	    
		placeHolder.append('<div class="collapse rawdataview" id=\"' + targetId + '\">' + result + '</div>');
		badge.removeClass('badge-secondary').addClass('badge-success').text('See Raw Data');
		badge.removeClass('need-moredata');
		
    	    },
    	    error: function (xhr, status, errorThrown) {
		alert("code:"+ xhr.status + "\n" + "message:" + xhr.responseText + "\n" + "error:" +error);
		badge.removeClass('badge-secondary').addClass('badge-primary').text(originalBadgeText);
    	    }
	});	
    });
    
    $('#master-playlist').on('click', '.variant-url', function() {
	if ($(this).children().attr('class').match('need-moredata')) {
	    m3u8Url = $(this).siblings('span').attr('data-item').trim();
	    targetId = $(this).siblings('span').attr('data-target-id').trim();	    
    	    appendMediaPlaylist(targetId, m3u8Url, $(this).attr('href').slice(1), $(this).parent(), $(this).children(), "Media Playlist");
	}	
    });

    $('#master-playlist').on('click', '.alternative-url', function() {	
	if ($(this).children().attr('class').match('need-moredata')) {
	    m3u8Url = $(this).siblings('span').attr('data-item').trim();
	    targetId = $(this).siblings('span').attr('data-target-id').trim();
	    appendMediaPlaylist(targetId, m3u8Url, $(this).attr('href').slice(1), $(this).parent(), $(this).children(), "Alternative Media Playlist");	    
	}	
    });

    $('#btn-clear').click(function() {
	$('#master-playlist').empty();
    });

    $('#btn-expand-all').click(function() {
	$('.collapse').collapse('show');
    });

    $('#btn-collapse-all').click(function() {
	$('.collapse').collapse('hide');
    });
});

function readM3u8(m3u8Url, placeholder) {
    placeholder.html('<h3>Loading Master Playlist...</h3>');
    $.ajax({
    	url: serverUrl,
    	type: 'POST',
    	data: {
    	    m3u8Url: m3u8Url
    	},
    	dataType: 'text',
    	success: function (result) {	    
	    placeholder.html(result);
    	},
    	error: function (xhr, status, errorThrown) {
    	    alert('failed to load : ' + m3u8Url);
	    placeholder.html('<h3>Failed to load</h3>');
    	}
    });
}

function appendMediaPlaylist(variantInfo, m3u8Url, elementId, element, badge, badgeText) {
    originalBadgeText = badge.text();
    badge.text("loading...");    
    $.ajax({
    	url: serverUrl,
    	type: 'POST',
    	data: {
    	    m3u8Url: m3u8Url,
	    variantInfo: variantInfo
    	},
    	dataType: 'text',
    	success: function (result) {	    
	    element.append('<ul class="collapse" id=\"' + elementId + '\">' + result + '</ul>');
	    badge.removeClass('badge-secondary').addClass('badge-success').text(badgeText);
	    badge.removeClass('need-moredata');
    	},
    	error: function (xhr, status, errorThrown) {
	    badge.text(originalBadgeText);
	    alert(errorThrown);
    	}
    });
}
