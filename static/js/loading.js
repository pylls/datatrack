function startLoading() {
  $("#main").hide();
  $("#loading").show();
}

function doneLoading() {
  $("#loading").fadeOut();
  $("#main").fadeIn();
}
