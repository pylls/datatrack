mapDT = 0
showMapCoords = true

jQuery(document).ready(function() {
		$("input[name='DTmapMode']").change(function(f){
			 showMapCoords = this.value == "coordinates";
			 updateMap();
	});

	mapDT = L.map('map', {worldCopyJump: true}).setView([59.402761, 13.514242], 13); //just define the default view manually
	L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
		attribution: '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors'
	}).addTo(mapDT);
  document.getElementById('map').style.display = 'block';
  mapDT.invalidateSize();

	mapDT.on('popupopen', function(e) {
		setDisclosureDetails("#" + e.popup._contentNode.childNodes[0].id)
	})

	mapDT.on('moveend', function(e) {
		updateMap()
	})

  updateMap()
});

function updateMap() {
	if (showMapCoords) {
		updateMapCoords()
	} else {
		updateMapLines()
	}
}

markers = new L.FeatureGroup()

function updateMapLines() {
	bounds = mapDT.getBounds()
	neLat = bounds.getNorthEast().lat
	neLng = bounds.getNorthEast().lng
	swLat = bounds.getSouthWest().lat
	swLng = bounds.getSouthWest().lng
	clearMapLines();
	markers.clearLayers();

	  $.get("/v1/coordinate/area/" + neLat.toString() + "/" + neLng.toString() + "/" + swLat.toString() + "/" + swLng.toString() + "/chronological", function(coordinates) {
			var points = new Array();
			for (i = 0; i < coordinates.length; i++) {
				if (coordinates[i].Prev.ID != "") {
					points.push(L.latLng(coordinates[i].Prev.Latitude, coordinates[i].Prev.Longitude));
				}

				points.push(L.latLng(coordinates[i].Latitude, coordinates[i].Longitude));

				if (coordinates[i].Next.ID != "") {
					points.push(L.latLng(coordinates[i].Next.Latitude, coordinates[i].Next.Longitude));
					setMapLine(points);
					points = new Array();
				}
			}
			setMapLine(points);
		})
}

function setMapLine(points) {
	var polyline = L.polyline(points, {});
	var decorator = L.polylineDecorator(polyline, {
		patterns: [
			{ offset: 0, repeat: 10, symbol: L.Symbol.dash({pixelSize: 5, pathOptions: {color: '#000', weight: 2, opacity: 0.5}}) },
			{offset: 0, repeat: '40px', symbol: L.Symbol.arrowHead({pixelSize: 5, pathOptions: {color: '#FFD600', weight: 6, opacity: 1}})}
		]}).addTo(mapDT);

	// var polyline = L.polyline(points, {color: "#FF0863", opacity: 0.75, dashArray: "5, 5, 1, 5", lineJoin: "miter", clickable: false});
	// markers.addLayer(polyline);
	// mapDT.addLayer(markers);
}

function clearMapLines() {
	mapDT.removeLayer(markers);

	for(i in mapDT._layers) {
		if(mapDT._layers[i]._path != undefined) {
			mapDT.removeLayer(mapDT._layers[i]);
		}
	}
}

function updateMapCoords() {
	bounds = mapDT.getBounds()
	neLat = bounds.getNorthEast().lat
	neLng = bounds.getNorthEast().lng
	swLat = bounds.getSouthWest().lat
	swLng = bounds.getSouthWest().lng

	mapDT.removeLayer(markers);
	markers.clearLayers();
	clearMapLines();

  // call for coordinates
  $.get("/v1/coordinate/area/" + neLat.toString() + "/" + neLng.toString() + "/" + swLat.toString() + "/" + swLng.toString(), function(coordinates) {

		for (i = 0; i < coordinates.length; i++) {
			var marker = L.marker([coordinates[i].Latitude, coordinates[i].Longitude]).bindPopup('<div class="dtpopup" id="popup' + coordinates[i].DisclosureID + '"><div class="dtpoptime">' + coordinates[i].Timestamp + '</div>');
			markers.addLayer(marker);
		}
		mapDT.addLayer(markers)
  })
}

function setDisclosureDetails(id) {
	did = id.substring(6)
	// basic popup structure
	$(id).append('<div class="dtpopcord dtpopstyle" /><div class="dtpopvel dtpopstyle" /><div class="dtpophead dtpopstyle" /><div class="dtpopacc dtpopstyle" /><div class="dtpopalt dtpopstyle" /><div class="dtpopact" />')

	// set time
	var timestamp = new Date(parseInt($(id + " > .dtpoptime").html(), 10));
	var dateStamp = timestamp.toUTCString().replace("GMT", "");
	var day = dateStamp.substring(0, 3);
	var date = dateStamp.substring(5,16);
	var time = dateStamp.substring(17,25);
	$(id + " > .dtpoptime").html(day + ', ' + date + ', ' + time);

	// read attributes
	$.get("/v1/disclosure/" + did + "/attribute", function(attributeIDs) {
		for(i = 0; i < attributeIDs.length; i++) {
			$.get("/v1/attribute/" + attributeIDs[i], function(attribute) {
				switch (attribute.Name) {
					case "Coordinates":
						$(id + " > .dtpopcord").html('<i class="fa fa-location-arrow fa-2x"></i>' + attribute.Value)
						break;
					case "Accuracy":
						$(id + " > .dtpopacc").html('<i class="fa fa-bullseye fa-2x"></i> Accuracy <br><span class="value">' + attribute.Value + '</span>')
						break;
					case "Altitude":
						$(id + " > .dtpopalt").html('<i class="fa fa-signal fa-2x"></i> Altitude <br><span class="value">' + attribute.Value + '</span>')
						break;
					case "Velocity":
						$(id + " > .dtpopvel").html('<i class="fa fa-tachometer fa-2x"></i> Velocity <br><span class="value">' + attribute.Value + '</span>')
						break;
					case "Heading":
						$(id + " > .dtpophead").html('<i class="fa fa-arrows-alt fa-2x"></i> Heading <br><span class="value">' + attribute.Value + '</span>')
						break;
				}
			})
		}
	})

	// read implicit attributes
	$.get("/v1/disclosure/" + did + "/implicit", function(disclosureIDs) {
		if(disclosureIDs.length != 0) {
				$('.dtpopact').before('<div class="activityheader">Activity(Confidence in %)</div>');
		}
		for(i = 0; i < disclosureIDs.length; i++) {
			$.get("/v1/disclosure/" + disclosureIDs[i] + "/attribute", function(attributeIDs) {
					for(j = 0; j < attributeIDs.length; j++) {
						$.get("/v1/attribute/" + attributeIDs[j], function(attribute) {
							activity = JSON.parse(attribute.Value);
							switch (activity.type) {
								case "tilting":
									$(id + " > .dtpopact").append('<div class="dtpopstyle dtpopstyleact"><img class="style" src="/img/shake.png"> '+ activity.type + '<br><span class="value">' + activity.confidence + '</span></div>')
									break;
								case "still":
									$(id + " > .dtpopact").append('<div class="dtpopstyle dtpopstyleact"><i class="fa fa-male fa-2x"></i> '+ activity.type + '<br><span class="value">' + activity.confidence + '</span></div>')
									break;
								case "walking":
									$(id + " > .dtpopact").append('<div class="dtpopstyle dtpopstyleact"><img class="style" src="/img/walk.png"> ' + activity.type + '<br><span class="value">'+ activity.confidence+ '</span></div>')
									break;
								case "onBicycle":
									$(id + " > .dtpopact").append('<div class="dtpopstyle dtpopstyleact"><i class="fa fa-bicycle fa-2x"></i> ' + activity.type + '<br><span class="value">'+ activity.confidence+ '</span></div>')
									break;
								case "inVehicle":
									$(id + " > .dtpopact").append('<div class="dtpopstyle dtpopstyleact"><i class="fa fa-car fa-2x"></i> ' + activity.type + '<br><span class="value">'+ activity.confidence+ '</span></div>')
									break;
								case "onFoot":
									$(id + " > .dtpopact").append('<div class="dtpopstyle dtpopstyleact"><img class="style" src="/img/foot.png"> ' + activity.type + '<br><span class="value">'+ activity.confidence+ '</span></div>')
									break;
								default:
									$(id + " > .dtpopact").append('<div class="dtpopstyle dtpopstyleact"><img class="style" src="/img/question.png"> '+ activity.type + '<br><span class="value">' + activity.confidence+ '</span></div>')
									break;
							}
						})
					}
			})
		}
	})
}
