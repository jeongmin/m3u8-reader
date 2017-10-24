$(document).ready(function() {
    $(function () {
	$('[data-toggle="tooltip"]').tooltip()
    })

    $('#read').click(function() {
	m3u8Url = $('#m3u8Url').val().trim();
	readM3u8(m3u8Url, '#master-playlist');
    });
    
    $('#master-playlist').on('click', '.variant-url', function() {	
	if ($(this).children().attr('class').match('badge-secondary')) {
	    m3u8Url = $(this).siblings('span').attr('data-item').trim();
	    targetId = $(this).siblings('span').attr('data-target-id').trim();
    	    appendMediaPlaylist(targetId, m3u8Url, $(this).attr('href').slice(1), $(this).parent(), $(this).children(), "Media Playlist");
	}	
    });

    $('#master-playlist').on('click', '.alternative-url', function() {	
	if ($(this).children().attr('class').match('badge-secondary')) {
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

function readM3u8(m3u8Url, selector) {
    $.ajax({
    	url: 'http://220.230.118.50:8080',
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

function appendMediaPlaylist(variantInfo, m3u8Url, elementId, element, badge, badgeText) {
    originalBadgeText = badge.text();
    badge.text("loading...");
    $.ajax({
    	url: 'http://220.230.118.50:8080',
    	type: 'POST',
    	data: {
    	    m3u8Url: m3u8Url,
	    variantInfo: variantInfo
    	},
    	dataType: 'text',
    	success: function (result) {	    
	    element.append('<ul class="collapse" id=\"' + elementId + '\">' + result + '</ul>');
	    badge.removeClass('badge-secondary').addClass('badge-info').text(badgeText);
    	},
    	error: function (xhr, status, errorThrown) {
	    badge.text(originalBadgeText);
	    alert(errorThrown);
    	}
    });
}
