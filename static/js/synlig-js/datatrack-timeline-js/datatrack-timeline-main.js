//make disclosureIDs global
disclosureIDs = 0;
shownDisclosureIDs = 0;
disclosuresPerLoad = 1;

$(document).ready(function() {
	
	//jump to top 
	 var offset = 220;
	    var duration = 500;
	    jQuery(window).scroll(function() {
	        if (jQuery(this).scrollTop() > offset) {
	            jQuery('.back-to-top').fadeIn(duration);
	        } else {
	            jQuery('.back-to-top').fadeOut(duration);
	        }
	    });
	    
	    jQuery('.back-to-top').click(function(event) {
	        event.preventDefault();
	        jQuery('html, body').animate({scrollTop: 0}, duration);
	        return false;
	    })
	
	// get all disclosureIDs
	$.get("/v1/disclosure/chronological/reverse", function(d) {
		disclosureIDs = d
		load(5);
	});


	//$('[data-toggle="tooltip"]').tooltip();
	$("[rel=tooltip]").tooltip();


	$("#cd-timeline").infinitus({
		trigger : function(done) {
			load(disclosuresPerLoad);
			setTimeout(done, disclosuresPerLoad*10);
		}
	});
});


function load(count) {
	if (shownDisclosureIDs + count > disclosureIDs.length) {
		count = disclosureIDs.length - shownDisclosureIDs
	}

	// create placeholder divs
	for (i = shownDisclosureIDs; i < count+shownDisclosureIDs; i++) {
		$("#cd-timeline").append("<div id=\"timeline-item-" + disclosureIDs[i] + "\"></div>");
	}

	// fill each div
	for (i = shownDisclosureIDs; i < count+shownDisclosureIDs; i++) {
		addDisclosure(disclosureIDs[i]);
	}
	shownDisclosureIDs += count;
}

function addDisclosure(id) {
	$.get("/v1/disclosure/" + id, function(disclosure) {
		// core part of disclosure
		$(document.getElementById("timeline-item-" + id)).addClass("cd-timeline-block").html(getDisclosureMarkup(disclosure));

		// attributes
		addAttributes(id, "attribute-space-" + id);

		// downstream attributes
		$.get("/v1/disclosure/" + id + "/implicit", function(dIDs) {
			if (dIDs.length > 0) {
				$("#attribute-space-derived-" + id).append('<h4 class="timeline-heading"><i class="fa fa-question-circle hoverchange" data-toggle="tooltip" title="Information the service extracted or infered from the information you shared"></i> Derived</h4>');
			}
			for (i = 0; i < dIDs.length; i++) {
				addDerivedDisclosure(dIDs[i], id)
			}
		});
	});
}

function addDerivedDisclosure(id, originDisclosureID) {
	$.get("/v1/disclosure/" + id, function(disclosure) {
		$("#attribute-space-derived-" + originDisclosureID).append(getDownstreamMarkup(disclosure));

		addAttributes(id, "downstream-attribute-space-" + id);
	});
}

function addAttributes(id, container) {
	$.get("/v1/disclosure/" + id + "/attribute", function(attributeIDs) {
		// add attributes
		for(i = 0; i < attributeIDs.length; i++) {
			$("#" + container).append('<div id="'+ container + id + i + '" />');
			addAttribute(attributeIDs[i], "#" + container + id + i);
		}
	});
}

function addAttribute(attributeID, id) {
	$.get("/v1/attribute/" + attributeID, function(attribute) {
		$(id).append(getAttributeMarkup(attribute));
	});
}

function getDownstreamMarkup(disclosure) {
	var timestamp = new Date(parseInt(disclosure.Timestamp, 10));
	var dateStamp = timestamp.toUTCString().replace("GMT", "");
	var day = dateStamp.substring(0, 3);
	var date = dateStamp.substring(5,16);
	var time = dateStamp.substring(17,25);

	return '<div class="timeline-derived-disclosure">' + day + ', ' + date + ', ' + time + '<br /><div id="downstream-attribute-space-' + disclosure.ID + '" />'
	+ '</div>';
}

function getDisclosureMarkup(disclosure) {
	var organizationUpperCase = disclosure.Recipient;
	var timestamp = new Date(parseInt(disclosure.Timestamp, 10));
	var dateStamp = timestamp.toUTCString().replace("GMT", "");;
	var day = dateStamp.substring(0, 3);
	var date = dateStamp.substring(5,16);
	var time = dateStamp.substring(17,25);
	var organization = organizationUpperCase.toLowerCase();

	return 	  '<div class="cd-timeline-img cd-priority-2 ">'  //cd-timleine-img is a css class for circles in the middle
	+			'<img class="clip-circle" src="../../img/iconsorganizations/'+organization+'.png" title="' + organizationUpperCase + '" />' //clip circle class is used for making images fit in the circle
	+		'</div>'
	+		'<div class="panel panel-default cd-timeline-content">'
	+			'<div class="timeline-panel-content">'
	+				'<div id="attribute-space-' + disclosure.ID + '" class="timeline-attributes-collection list-group list-group-noborder">' // where we put attributes
	+       '<h4 class="timeline-heading">'
	+				'<i  id="test" class="fa fa-question-circle hoverchange" rel="tooltip" data-toggle="tooltip" title="Information you or one of your devices shared with the service"></i>'
	+					' Disclosed by you'
	+		  '</h4>'
	+		'</div>'
	+				'<div id="attribute-space-derived-' + disclosure.ID + '" class="timeline-derived-attributes-collection list-group list-group-noborder" />' // derived attributes
	+			'</div>'
	+				'<span class="cd-date">'
	+               	'<p class="disclosure-day">'
	+						day + '&nbsp;'
	+					'</p>'
	+               	'<p class="disclosure-date">'
	+						date + '&nbsp;'
	+					'</p>'
	+               	'<p class="disclosure-time">'
	+						'' + time + '&nbsp;'
	+					'</p>'
	+ 				'</span>'
	+		'</div><!--cd-timeline-content-->'
	;
}

function getAttributeMarkup(attribute) {
	if (attribute.Name == "Activity") {
		activity = JSON.parse(attribute.Value);

		var activityvalue = "";

		switch (activity.type) {
		case "tilting":
			activityvalue = '<img class="style" src="/img/shake.png"> ';
			break;
		case "still":
			activityvalue = '<i class="fa fa-male"></i> ';
			break;
		case "walking":
			activityvalue = '<img class="style" src="/img/walk.png"> ';
			break;
		case "onBicycle":
			activityvalue = '<i class="fa fa-bicycle"></i> ';
			break;
		case "inVehicle":
			activityvalue = '<i class="fa fa-car"></i> ';
			break;
		case "onFoot":
			activityvalue = '<img class="style" src="/img/foot.png"> ';
			break;
		default:
			activityvalue ='<img class="style" src="/img/question.png"> ';
		break;
		}
		return	'<div class ="timeline-attribute">'
		+					'<div>'
		+						'<div class="attribute-name">'
		// +							'<span class="fa fa-'+ attribute.Type +'"> &nbsp&nbsp '
		+							activityvalue
		+								attribute.Name
		+							'</span>'
		+						'</div>'
		+						'<div class="attribute-value">'
		+							'<p class="disclosure-text-properties">'+ activity.type + ', ' + activity.confidence +'% confidence</p>'
		+						'</div>'
		+					'</div>';
	}
	return	'<div class ="timeline-attribute">'
	+					'<div>'
	+						'<div class="attribute-name">'
	+							'<span class="fa fa-'+ attribute.Type +'"> &nbsp&nbsp '
	+								attribute.Name
	+							'</span>'
	+						'</div>'
	+						'<div class="attribute-value">'
	+							'<p class="disclosure-text-properties">'+ attribute.Value +'</p>'
	+						'</div>'
	+					'</div>';
}
