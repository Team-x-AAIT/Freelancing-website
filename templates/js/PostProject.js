function setFiles() {
    var uploader = document.getElementById('file-uploader');
    uploader.click();
}

function previewFiles(uploader) {
    var element = document.getElementById("pdf-display");
    var elementOuter = document.getElementById("pdf-display-outer");
    elementOuter.style = "display:block !important";
    for (var i = 0; i < uploader.files.length; i++){
        if (uploader.files && uploader.files[i]) {
            
            var pdfFile = uploader.files[i];
            var name = pdfFile.name;

            var div = document.createElement("div");
            var para = document.createElement("p");
            var nodeText = document.createTextNode(name);
            var iconPdf = document.createElement("i");
            para.style = "font-size:13px; float:left; margin:0 5px 0";
            iconPdf.className = "fa fa-file-pdf-o";
            iconPdf.style = "font-size:30px;color:red; float:left";
    
            para.appendChild(nodeText);
            div.append(iconPdf);
            div.append(para);
            element.appendChild(div);
        }
    
    }
}

function removeFiles(){

    $("#pdf-display").empty();

    var elementOuter = document.getElementById("pdf-display-outer");
    var uploader = document.getElementById('file-uploader');
    elementOuter.style = "display:none !important";
    uploader.value = "";
}


function onPostProjectButtonClicked() {
   
    var title = document.getElementById('smallDescription').value;
    var category = document.getElementById('category').value;
    var subCategory = document.getElementById('subcategory').value;
    var description = document.getElementById('description').value;
    var details = document.getElementById('details').value;
    var workType = document.getElementById('workType').value;
    var budget = document.getElementById('budget').value;
    var uploader = document.getElementById('file-uploader');
    var alertMessage = document.getElementById('alert-message');
    var submit = document.getElementById('submit');
    alertMessage.style = "color: red; font-size: small; clear: both;display: block !important";
    $("#alert-message").empty();
    var index;
    var message;

    if (title.length < 3) {
        message = "Title should be at least 3 characters!";
        alertMessage.appendChild(document.createTextNode(message));
        return;
    }

    if (description.length < 3) {
        message = "Description should be at least 3 characters!";
        alertMessage.appendChild(document.createTextNode(message));
        return;
    }

    if (details.length < 3) {
        message = "Details should be at least 3 characters!";
        alertMessage.appendChild(document.createTextNode(message));
        return;
    }
    
    if (category == 0) {
        message = "Please select category!";
        alertMessage.appendChild(document.createTextNode(message));
        return;
    }

    if (subCategory == 0) {
        message = "Please select subcategory!";
        alertMessage.appendChild(document.createTextNode(message));
        return;
    }

    if ((workType == '1' || workType == '2') && budget == '') {
        message = "Please fill budget!";
        alertMessage.appendChild(document.createTextNode(message));
        return;
    }

    submit.click();



    // var data = new FormData();

    // data.append('title', title);
    // data.append('category', category);
    // data.append('subcategory', subCategory);
    // data.append('description', description);
    // data.append('details', details);
    // data.append('worktype', workType);
    // data.append('budget', budget);
    // data.append('files', uploader.files[0]);

    // $.ajax({
    //     url: "/Post_Project",
    //     type: "POST",
    //     data: data,
    //     processData: false,
    //     contentType: false,
    //     success: function (response) {
    //         switch (response) {
    //             case "title length too short": index = 1; message = "Title should be at least 3 characters!"; return;
    //             case "unknown category": index = 2; messae = "Please fill category!"; return;
    //             case "unknown subcategory": index = 3; messae = "Please fill subcategory!"; return;
    //             case "description length too short": index = 4; message = "Description should be at least 3 characters!"; return;
    //             case "details length too short": index = 5; message = "Details should be at least 3 characters!"; return;
    //             case "unknown worktype": index = 6; messae = "Please Fill work type"; return;

    //         }
    //         if (index > 0) {
    //             alertMessage.appendChild(document.createTextNode(message));
    //             return;
    //         }
    //         if (response == "okay") {
    //             window.location.href = "/Dashboard";
    //         }
    //     }
    // });
    
}