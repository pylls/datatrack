$(document).ready(function() {
	
	var ifUploaded = false;
	
	$("#start").show();

	$('#fileToUpload').change(function() {
		$('#toProcessing').prop('disabled', !($('#fileToUpload').val()));
	});

	$.get("/v1/disclosure/count", function(disclosure) {  //check if file has been uploaded
		if(disclosure != 0)
			ifUploaded  = true;
	});
	
	$("#back").click(function() {
		
		$("#import").toggleClass("hide showme");
		$("#uploadpart").toggleClass("hide showme");
		$("#processing-error").toggleClass("hide showme");

	});

	
	$("#toImport").click(function() {
		$("#start").hide();
		$("#import").toggleClass("hide showme");
		if(!ifUploaded)
		{
			$("#uploadpart").toggleClass("hide showme");

		}
		else 
		{
			var html ='<p>You have already uploaded your file.</p>' 
						+'</br>' 
						+'<button id="nextstep" type="button"	class="btn btn-primary btn-lg">Start Exploring</button>';	
			
			$("#textcomplete").html(html);
			$("#nextstep").click(function() {
				$("#uploadpart").toggleClass("hide showme");
				$("#import").toggleClass("hide showme");
				$("#done").toggleClass("hide showme");
			});
		}
	});

	$("#importform").submit(function(event) {
		event.stopPropagation();
		event.preventDefault();
		$("#uploadpart").toggleClass("hide showme");
		$("#import").toggleClass("hide showme");
		$("#processing").toggleClass("hide showme");

		var formData = new FormData($(this)[0]);

		$.ajax({
			url: "/v1/google",
			type: "POST",
			data: formData,
			success: function() {
				$("#processing").toggleClass("hide showme");
				$("#done").toggleClass("hide showme");
			},
			error: function() { 
				//alert("error"); 
				$("#processing").toggleClass("hide showme");

				$("#processing-error").toggleClass("hide showme");				
			},
			contentType: false,
			cache: false,
			processData: false
		});
	});

});

function importComplete() {
	alert("done");
}

function importError() {
	alert("ohnoe");
}
