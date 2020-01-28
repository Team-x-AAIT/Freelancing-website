window.addEventListener('load', onProfilePageLoad);

function setProfileImage() {
    var uploader = document.getElementById('image-uploader');
    uploader.click();
}

function previewProfileImage(uploader) {
    if (uploader.files && uploader.files[0]) {
        var imageFile = uploader.files[0];
        var reader = new FileReader();
        reader.onload = function (e) {
            $('#profile-pic').attr('src', e.target.result);
        };
        reader.readAsDataURL(imageFile);
    }
}

function setCV() {
    var uploader = document.getElementById('cv-uploader');
    uploader.click();
}

function previewCV(uploader) {
    if (uploader.files && uploader.files[0]) {
        var element = document.getElementById("pdf-display");

        if (element.childElementCount != 0) {
            element.removeChild(element.childNodes[1]);
        }

        var pdfFile = uploader.files[0];
        var name = pdfFile.name;

        var div = document.createElement("div");
        var para = document.createElement("p");
        var nodeText = document.createTextNode(name);
        var iconPdf = document.createElement("i");
        var iconx = document.createElement("i");
        para.style = "font-size:13px; float:left; margin:0 5px 0";
        iconPdf.className = "fa fa-file-pdf-o";
        iconPdf.style = "font-size:30px;color:red; float:left";
        iconx.className = "fas fa-times";
        iconx.style = "font-size: 10px;color:red; float:left; cursor: pointer";
        iconx.onclick = removeResource;

        para.appendChild(nodeText);
        div.append(iconPdf);
        div.append(para);
        div.append(iconx);


        element.appendChild(div);

    }
}

function onProfilePageLoad() {


    var csrfError = document.getElementById("csrf_error");
    var csrfErrorDiv = document.getElementById("csrf_error_div");
    if (csrfError.value != "") {
        csrfErrorDiv.style.display = "block";
    }

    var hiddenInput1 = document.getElementById("image-value");
    var profilePic = document.getElementById("profile-pic");
    var imageSrc = hiddenInput1.value;

    var hiddenInput2 = document.getElementById("cv-value");
    var pdfRepo = document.getElementById("pdf-repository");
    var filename = hiddenInput2.value;

    var hiddenInput3 = document.getElementById("gender-value");
    var male = document.getElementById("male");
    var female = document.getElementById("female");
    var other = document.getElementById("other");
    var genderValue = hiddenInput3.value;

    var imageUploder = document.getElementById('image-uploader');
    var cvUploder = document.getElementById('cv-uploader');
    imageUploder.files[0] = "";
    cvUploder.files[0] = "";

    if (filename != "") {
        var div = document.createElement("a");
        var para = document.createElement("p");
        var nodeText = document.createTextNode(filename);
        var iconPdf = document.createElement("i");

        para.style = "font-size:13px; float:left; margin:0 5px 0";
        iconPdf.className = "fa fa-file-pdf-o";
        iconPdf.style = "font-size:30px;color:red; float:left";
        div.style = "overflow:hidden";
        div.href = "/asset_viewer?type=cv&asset=" + filename;
        div.target = "_blank";

        para.appendChild(nodeText);
        div.append(iconPdf);
        div.append(para);

        pdfRepo.appendChild(div);
        pdfRepo.style = "display:inline-block !important";
    }

    if (imageSrc != "") {
        profilePic.src = "/assets/profile_pic/" + imageSrc;
    }

    switch (genderValue) {
        case "M":
            male.checked = true;
            break;
        case "F":
            female.checked = true;
            break;
        case "O":
            other.checked = true;
            break;
    }

    spinner.style = "display: none !important";

}


function removeResource() {

    var uploader = document.getElementById('cv-uploader');
    uploader.value = "";
    $("#pdf-display").empty();

}

function statusClicked(statusBtn) {
    var hiddenStatus = document.getElementById('hiddenStatus');
    var statusBtn1 = document.getElementById('status-btn1');
    var statusBtn2 = document.getElementById('status-btn2');
    if (statusBtn.id == "status-btn1") {
        hiddenStatus.value = "1";
        statusBtn.style.color = "#86bce9";
        statusBtn2.style.color = "white";
    } else if (statusBtn.id == "status-btn2") {
        hiddenStatus.value = "0";
        statusBtn.style.color = "#86bce9";
        statusBtn1.style.color = "#da571f";
    }
}

function onSaveButtonClicked() {

    var spinner = document.getElementById('spinner');
    var firstName = document.getElementById('firstname').value;
    var lastName = document.getElementById('lastname').value;
    var jobTitle = document.getElementById('jobTitle').value;
    var email = document.getElementById('email').value;
    var phoneNumber = document.getElementById('phonenumber').value;
    var aboutYou = document.getElementById('bio').value;
    var country = document.getElementById('country').value;
    var city = document.getElementById('city').value;
    var gender = $("input[name='gender']:checked").val();
    // var prefer = 'yes';
    var imageUploder = document.getElementById('image-uploader');
    var cvUploder = document.getElementById('cv-uploader');

    spinner.style = "display: inline-block !important";

    // if ($('#prefer').prop("checked") == false) {
    //     prefer = "no";
    // }

    if (firstName == "") {
        message = "Please fill first name";
        errorCheck(message);
        return;
    }
    if (lastName == "") {
        message = "Please fill last name";
        errorCheck(message);
        return;
    }
    if (email == "" || !validateEmail(email)) {
        message = "Please use valid email";
        errorCheck(message);
        return;
    }
    if (!validatePhoneNumber(phoneNumber)) {
        message = "Please use valid phone number";
        errorCheck(message);
        return;
    }

    var profilePic = imageUploder.files[0];
    var cv = cvUploder.files[0];
    var data = new FormData();

    data.append('firstname', firstName);
    data.append('lastname', lastName);
    data.append('phonenumber', phoneNumber);
    data.append('email', email);
    data.append('jobTitle', jobTitle);
    data.append('bio', aboutYou);
    // data.append('prefer', prefer);
    data.append('country', country);
    data.append('city', city);
    data.append('gender', gender);
    data.append('profilePic', profilePic);
    data.append('cv', cv);
    // data.append('saveFlag', 'true');

    $.ajax({
        url: "/EditProfile/Update",
        type: "POST",
        data: data,
        processData: false,
        contentType: false,
        success: function (msg) {
            changeProfile(msg);
        }
    });
}

function changeProfile(msg) {

    if (msg == 'Okay') {
        window.location.replace("/Dashboard");

    } else {
        switch (msg) {
            case "invalid firstname":
                message = "Please fill first name";
                errorCheck(index, message);
                return;

            case "invalid lastname":
                message = "Please fill last name";
                errorCheck(message);
                return;

            case "invalid email":
                message = "Please use valid email";
                errorCheck(message);
                return;

            case "invalid phoneNumber":
                message = "Please use valid phone number";
                errorCheck(message);
                return;
            default:
                message = "Oops! something went wrong";
        }
    }
}

function errorCheck(erroMessage) {
    var element = document.getElementById("error-message");
    var nodeText = document.createTextNode(erroMessage);
    element.appendChild(nodeText);
    spinner.style = "display: none !important";
}


function validateEmail(email) {
    var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
}

function validatePhoneNumber(phoneNumber) {
    phoneNumber = phoneNumber.trim();
    if (phoneNumber.length != 9) {
        return false;
    } else if (phoneNumber.charAt(0) != 9) {
        return false;
    }
    return true;
}

function onFocusInInput() {
    $("#error-message").empty();
}