window.addEventListener('load', onViewProjectLoad);

function onViewProjectLoad(){
    var element = document.getElementById("hiddenWorkType");
    var workType = element.value;
    switch (workType) {
        case "1":
            workType = "Fixed";
            break;
        case "2":
            workType = "Perhour";
            break;
        case "3":
            workType = "Negotaiable";
            break;
        default:
            workType = "Unkown";
    }
    var pWorkType = document.getElementById("article-project-workType");
    pWorkType.appendChild(document.createTextNode(workType));

    var profilePic = document.getElementById("project-owner-img");
    var hiddenOwnerImg = document.getElementById("hiddenOwnerImg");
    var imageSrc = hiddenOwnerImg.value;

    if (imageSrc != "") {
        profilePic.src = "/assets/profile_pic/" + imageSrc;
    }

}


function onApplyClicked(){
    var pid = document.getElementById("pid").value;
    var proposal = document.getElementById("message").value;
    var alertBox = document.getElementById("alert-box");
    var alertMessage = document.getElementById("alert-message");

    var data = new FormData();
    data.append('pid', pid);
    data.append('proposal', proposal);

    $.ajax({
        url: "/ApplyFor_Project",
        type: "POST",
        data: data,
        processData: false,
        contentType: false,
        success: function (msg) {

            if (msg == "okay"){
                window.location.href = "/Dashboard";
            }else{
                alertBox.style = "display: block !important;";
                alertMessage.appendChild(document.createTextNode(msg + "!"));
                alertMessage.style = "font-weight:bold;";
            }
        }
    });
}