window.addEventListener('load', onPageLoad);

function onPageLoad() {

  var csrfError = document.getElementById("csrf_error");
  var csrfErrorDiv = document.getElementById("csrf_error_div");
  if (csrfError.value != "") {
    csrfErrorDiv.style.display = "block";
  }
}

function openNav() {
  document.getElementById("mySidenav").style.width = "250px";
  document.getElementById("main").style.marginLeft = "250px";
}

function closeNav() {
  document.getElementById("mySidenav").style.width = "0";
  document.getElementById("main").style.marginLeft = "0";
}

var email = "";
var password = "";
var rememberMe = "";


function loginButtonClicked() {
  document.getElementById("alert").innerHTML = "";
  document.getElementById("buttonSpinner").className = "";
  if (email == "" || password == "") {
    email = document.getElementsByTagName("input")[0].value;
    password = document.getElementsByTagName("input")[1].value;
    rememberMe = document.getElementsByTagName("input")[2].checked;
  }

  if (email == "" || password == "") {
    return;
  }

  document.getElementById("buttonSpinner").className =
    "spinner-border spinner-border-sm";
  var request;
  if (window.XMLHttpRequest) {
    request = new XMLHttpRequest();
  } else {
    request = new ActiveXObject("Microsoft.XMLHTTP");
  }

  request.onreadystatechange = function() {
    if (request.readyState == 4 && request.status == 200) {
      var response = request.responseText;

      if (response == "invalid") {
        document.getElementById("alert").innerHTML =
          "Invalid username or password!";
        document.getElementById("buttonSpinner").className = "d-none";
      }

      if (response == "okay") {
        window.location.replace("/Dashboard");
        document.getElementById("buttonSpinner").className = "d-none";
      }
    }
  };
  request.open("POST", "/Login", true);
  request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  request.send(
    "email=" + email + "&password=" + password + "&rememberMe=" + rememberMe
  );
}

function LoginWGoogle(){
  window.location.href = "/GoogleLogin";
}

function LoginWFacebook() {
  window.location.href = "/FacebookLogin";
}

function LoginWLinkedIn() {
  window.location.href = "/LinkedInLogin";
}