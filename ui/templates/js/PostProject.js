window.addEventListener('load', onPageLoad);


function setFiles() {
    var uploader = document.getElementById('file-uploader');
    uploader.click();
}

function onPageLoad(){

    var csrfError = document.getElementById("csrf_error");
    var csrfErrorDiv = document.getElementById("csrf_error_div");
    if (csrfError.value != "") {
        csrfErrorDiv.style.display = "block";
    }

    var hiddenCategory = document.getElementById("hiddenCategory");
    var hiddenWorkType = document.getElementById("hiddenWorkType");

    var category = document.getElementById("category");
    
    var workType = document.getElementById("workType");

    if (hiddenCategory.value != ""){
        childs = category.childNodes;
        for (var i = 0; i < childs.length; i++){
            if (childs[i].value == hiddenCategory.value){
                category.selectedIndex = i - 2;
                category.onchange();
                break;
            }
        }
        
    }
    if (hiddenWorkType.value != "") {
        childs = workType.childNodes;
        for (var k = 0; k < childs.length; k++) {
            if (childs[k].value == hiddenWorkType.value) {
                workType.selectedIndex = k - 1;
                break;
            }
        }
    }


}

function previewFiles(uploader) {
    var element = document.getElementById("pdf-display");
    var elementOuter = document.getElementById("pdf-display-outer");
    elementOuter.style = "display:block !important";
    for (var i = 0; i < uploader.files.length; i++) {
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

function removeFiles() {

    $("#pdf-display").empty();

    var elementOuter = document.getElementById("pdf-display-outer");
    var uploader = document.getElementById('file-uploader');
    elementOuter.style = "display:none !important";
    uploader.value = "";
}


function onCategoryChange() {

    var hiddenSubcategory = document.getElementById("hiddenSubcategory");
    var subCategory = document.getElementById("subcategory");
    subCategory.innerHTML = "";
    subCategory.disabled = true;

    var optionN = document.createElement("option");
    optionN.value = "none";
    optionN.text = "Select";
    subCategory.appendChild(optionN);
    
    var data = new FormData();

    var categoryValue = $("#category").children("option:selected").val();
    data.append('category', categoryValue);

    $.ajax({
        url: "/SubCategories",
        type: "POST",
        data: data,
        processData: false,
        contentType: false,
        success: function (response) {
            var tempArray = JSON.parse(response);
            if (tempArray != null && tempArray.length > 0) {
                for (var i = 0; i < tempArray.length; i++) {
                    var option = document.createElement("option");
                    option.value = tempArray[i];
                    option.text = tempArray[i];
                    subCategory.appendChild(option);
                }
                subCategory.disabled = false;

                if (hiddenSubcategory.value != "") {
                    childs = subCategory.childNodes;
                    for (var j = 0; j < childs.length; j++) {
                        if (childs[j].value == hiddenSubcategory.value) {
                            subCategory.selectedIndex = j;
                            break;
                        }
                    }
                    return;
                }
            }
        }
    });

}