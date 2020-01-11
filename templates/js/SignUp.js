function signUpBtnClicked() {
  document.getElementById("alert").innerHTML = "";
  document.getElementById("buttonSpinner").className = "";

  firstname = document.getElementById("firstname").value;
  lastname = document.getElementById("lastname").value;
  email = document.getElementById("email").value;
  password = document.getElementById("password").value;
  cPassword = document.getElementById("confirmPassword").value;

  if (
    email == "" ||
    password == "" ||
    firstname == "" ||
    lastname == "" ||
    cPassword == ""
  ) {
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

      if (response == "okay") {

        window.location.href = "/Check_Your_Email";
        document.getElementById("buttonSpinner").className = "d-none";
      }
      else {
        document.getElementById("alert").innerHTML = response + "!";
        document.getElementById("buttonSpinner").className = "d-none";
      }
    }
  };
  request.open("POST", "/Register", true);
  request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  request.send(
    "email=" +
      email +
      "&password=" +
      password +
      "&firstname=" +
      firstname +
      "&lastname=" +
      lastname +
      "&confirmPassword=" +
      cPassword +
      "&thirdParty=" +
      "false"
  );
}
