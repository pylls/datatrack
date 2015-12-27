function onLinkedInLoad() {
	IN.Event.on(IN, "auth", getProfileData);

}
function onLinkedInLogout() {
	document.getElementById("logout").innerHTML = "you are logged out from Linkedin";
}
// Handle the successful return from the API call
function onSuccess(data) {
	console.log(data);
}

// Handle an error response from the API call
function onError(error) {
	console.log(error);
}

// Use the API call wrapper to request the member's basic profile data
function getProfileData() {	
	IN.API.Raw("/people/~").result(onSuccess).error(onError);
	IN.API.Raw('/people/~:(id,num-connections,picture-url)')
	.result(function(value) {
		console.log(value);
	})
	.error(function(error) {
		console.log(JSON.stringify(error));
	});

	IN.API.Raw("/people/~:(firstName,lastName,positions,skills,languages:(id,language,proficiency),educations)")
	.result(onSuccess).error(onError);
}

window.fbAsyncInit = function() {
	FB.init({
		appId      : '1123138574379887',
		xfbml      : true,
		version    : 'v2.3'
	});

	FB.getLoginStatus(function(response)
			{
		if (response.status === 'connected')
		{
			// the user is logged in and has authenticated your
			// app, and response.authResponse supplies
			// the user's ID, a valid access token, a signed
			// request, and the time the access token 
			// and signed request each expire
			var uid = response.authResponse.userID;
			var accessToken = response.authResponse.accessToken;
		//	$("#facebook").html('<i class="fa fa-check-square-o"></i> You are connected to Facebook');
			getinformation(uid,accessToken);
		}
		else if (response.status === 'not_authorized')
		{
			// the user is logged in to Facebook, 
			// but has not authenticated your app
			 //$("#facebook").html('<i class="fa fa-times"></i> You are disconnected from Facebook');

		}
		else
		{
			// $("#facebook").html('<i class="fa fa-times"></i> You are disconnected from Facebook');

			//the user is not logged in to the facebook
		}
			});
};

function checkLoginState() {
	FB.getLoginStatus(function(response) {

		if (response.status === 'connected')
		{
			// the user is logged in and has authenticated your
			// app, and response.authResponse supplies
			// the user's ID, a valid access token, a signed
			// request, and the time the access token 
			// and signed request each expire
			var uid = response.authResponse.userID;
			var accessToken = response.authResponse.accessToken;
			//$("#facebook").html('<i class="fa fa-check-square-o"></i> You are connected to Facebook');
			getinformation(uid,accessToken);

		}
		//else  $("#facebook").html('<i class="fa fa-times"></i> You are disconnected from Facebook');


	});
}
function getinformation(uid,accessToken)
{
	var uri = "/" + uid + "?access_token=" + accessToken;
	var uri2 = "/" + uid + "?access_token=" + accessToken+"?fields=friends";
	FB.api(
			uri,
			function (response) {
				if (response && !response.error) {
					console.log(response);

				}
			}
	);

	FB.api(
			"me?fields=friends",
			function (response) {
				if (response && !response.error) {
					console.log(response);

				}
			}
	);



}


function testAPI() {
	console.log('Welcome!  Fetching your information.... ');
	FB.api('/me', function(response) {
		console.log('Successful login for: ' + response.name);
		document.getElementById('status').innerHTML =
			'Thanks for logging in, ' + response.name + '!';
	});
}


////for linkedin

// Setup an event listener to make an API call once auth is complete
function insertdata()
{
	
	 $.post("/v1/import", { Organization: "CardioMon", owner : "vasilis"} , function(d)
			 {
		 		console.log(d);
			 }
	 ).success(function() {
			alert("Data has been inserted");
			$("#inserteddata").html("You have inserted CardioMon data");
		}).error(function() {
			alert("error while inserting data");
		});


}