window.addEventListener('load', onDashboardLoad);
var count = 1;

function onDashboardLoad(){
    $.ajax({
        url: "/Get_Match_Tag",
        type: "POST",
        processData: false,
        contentType: false,
        success: function (msg) {

            var listOfMatchTag = JSON.parse(msg);
            if (listOfMatchTag != null){
                listOfMatchTag.forEach(element => {
                    viewMatchTag(element);
                });
            }
        }
    });

    GetProjects();
}

function addMatchTag(){
    if (count < 4){
        var element = document.getElementById("match-container");
        var category = document.getElementById("category");
        var subcategory = document.getElementById("subcategory");
        var worktype = document.getElementById("worktype");
        var div = document.createElement("div");
        var hiddenInput1 = document.createElement("input");
        var hiddenInput2 = document.createElement("input");
        var hiddenInput3 = document.createElement("input");
        var nodeText = document.createTextNode("Match -" + count);
        var iconx = document.createElement("i");

        hiddenInput1.value = category.value;
        hiddenInput2.value = subcategory.value;
        hiddenInput3.value = worktype.value;

        div.className = "text-center p-sm-0";
        div.style = "cursor: pointer";
        iconx.className = "fas fa-times";
        iconx.style = "color:red; float: right; cursor: pointer";
        hiddenInput1.hidden = true;
        hiddenInput2.hidden = true;
        hiddenInput3.hidden = true;

        iconx.onclick = function () {
            div.remove();
            count--;
            $.ajax({
                url: "/Remove_Match_Tag",
                type: "POST",
                data: data,
                processData: false,
                contentType: false,
                success: function (msg){
                   GetProjects();
                }
            });

        };

        div.ondblclick = function(){
            var toogleButton = document.getElementById("match-adder");

            category.value = hiddenInput1.value;
            subcategory.value = hiddenInput2.value;
            worktype.value = hiddenInput3.value;

            toogleButton.click();    

        };

        var data = new FormData();

        data.append('category', hiddenInput1.value);
        data.append('subcategory', hiddenInput2.value);
        data.append('worktype', hiddenInput3.value);


        $.ajax({
            url: "/Add_Match_Tag",
            type: "POST",
            data: data,
            processData: false,
            contentType: false,
            success: function (msg) {
                if (msg == "okay"){
                    div.appendChild(nodeText);
                    div.append(iconx);
                    element.append(div);

                    count++;
                    GetProjects();
                }
            }
        });
       
    }
    
}

function viewMatchTag(matchTag){

    var element = document.getElementById("match-container");

    var div = document.createElement("div");
    var hiddenInput1 = document.createElement("input");
    var hiddenInput2 = document.createElement("input");
    var hiddenInput3 = document.createElement("input");
    var nodeText = document.createTextNode("Match -" + count);
    var iconx = document.createElement("i");

    hiddenInput1.value = matchTag.Category;
    hiddenInput2.value = matchTag.Subcategory;
    switch (matchTag.WorkType) {
        case 1:
            hiddenInput3.value = "Fixed";
            break;
        case 2:
            hiddenInput3.value = "Perhour";
            break;
        case 3:
            hiddenInput3.value = "Negotiable";
            break;
        case 4:
            hiddenInput3.value = "";
            break;
    }

    div.className = "text-center p-sm-0";
    div.style = "cursor: pointer";
    iconx.className = "fas fa-times";
    iconx.style = "color:red; float: right; cursor: pointer";
    hiddenInput1.hidden = true;
    hiddenInput2.hidden = true;
    hiddenInput3.hidden = true;

    var data = new FormData();

    data.append('category', hiddenInput1.value);
    data.append('subcategory', hiddenInput2.value);
    data.append('worktype', hiddenInput3.value);

    iconx.onclick = function () {
        div.remove();
        count--;
        $.ajax({
            url: "/Remove_Match_Tag",
            type: "POST",
            data: data,
            processData: false,
            contentType: false,
            success: function (msg){
                GetProjects();
            }
        });

    };

    div.ondblclick = function () {
        var toogleButton = document.getElementById("match-adder");

        category.value = hiddenInput1.value;
        subcategory.value = hiddenInput2.value;
        worktype.value = hiddenInput3.value;

        toogleButton.click();

    };

    div.appendChild(nodeText);
    div.append(iconx);
    element.append(div);

   count++;
}

function GetProjects(){

    var invisibleDiv = document.getElementById("invisible-div-2");
    var spinner = document.getElementById("spinner-freelancers");

    invisibleDiv.style = "display:none !important";
    spinner.style = "display: block !important; padding-top: 220px; background-color: #FCFAFA; height: 550px";
    $.ajax({
        url: "/Get_Projects",
        type: "POST",
        processData: false,
        contentType: false,
        success: function (msg) {

            projectList = JSON.parse(msg);
            if (projectList != null) {
                projectList.forEach(element => {
                    switch(element.Project.WorkType){
                        case 1:
                            element.Project.WorkType = "Fixed";
                            break;
                        case 2:
                            element.Project.WorkType = "Perhour";
                            break;
                        case 3:
                            element.Project.WorkType = "Negotiable";
                            
                    }
                    clearCards();
                    viewProjects(element);
                });
                spinner.style = "display:none !important";
                return;
            }
            clearCards();
            invisibleDiv.style = "display:block !important";
            spinner.style = "display:none !important";
            
        }
    });

}


function viewProjects(container){

    if (container.Project.Title == "" || container.Project.Title == null) {
        return;
    }

    var element = document.getElementById("freelancers-div");
    var div = document.createElement("div");
    var pTitleDiv = document.createElement("div");
    var imgDiv = document.createElement("div");  
    var image = document.createElement("img");
    var pTitle = document.createElement("p");
    var pDescription = document.createElement("p");
    var pDescriptionValue = document.createElement("p");
    var typeDiv = document.createElement("div");
    var maskedDiv = document.createElement("div");
    var unMaskedDiv = document.createElement("div");
    var pCategory = document.createElement("p");
    var pSubcategory = document.createElement("p");
    var insideTypeDiv = document.createElement("div");
    var pType = document.createElement("p");
    var pBudget = document.createElement("p");
    var pDate = document.createElement("p");
    var categorySpan = document.createElement("span");
    var subcategorySpan = document.createElement("span");
    var typeSpan = document.createElement("span");
    var budgetSpan = document.createElement("span");
    var typeSpanValue = document.createElement("span");
    var budgetSpanValue = document.createElement("span");
    var budgetSpanETB = document.createElement("span");

    var titleTextNode = document.createTextNode(container.Project.Title);
    var descriptionValueTextNode = document.createTextNode(container.Project.Description);
    var categoryTextNode = document.createTextNode(container.Project.Category);
    var subcategoryTextNode = document.createTextNode(container.Project.Subcategory);
    var typeTextNode = document.createTextNode(container.Project.WorkType);
    var budgetTextNode = document.createTextNode(container.Project.Budget);
    var date = new Date(container.Project.CreatedAt);
    var dateTextNode = document.createTextNode(date.toDateString());



    element.append(div);
    pTitleDiv.append(pTitle);
    imgDiv.append(image);
    pTitleDiv.append(imgDiv);
    div.append(pTitleDiv);
    div.append(pDescription);
    div.append(pDescriptionValue);
    div.append(typeDiv);
    div.append(pDate);

    typeDiv.append(maskedDiv);
    typeDiv.append(unMaskedDiv);

    maskedDiv.append(pCategory);
    maskedDiv.append(pSubcategory);
    pCategory.append(categorySpan);
    pSubcategory.append(subcategorySpan);

    unMaskedDiv.append(insideTypeDiv);
    insideTypeDiv.append(pType);
    insideTypeDiv.append(pBudget);
    pType.append(typeSpan);
    pType.append(typeSpanValue);
    pBudget.append(budgetSpan);
    pBudget.append(budgetSpanValue);
    pBudget.append(budgetSpanETB);


    pTitle.appendChild(titleTextNode);
    pDescription.appendChild(document.createTextNode("Description:"));
    pDescriptionValue.appendChild(descriptionValueTextNode);
    categorySpan.appendChild(categoryTextNode);
    subcategorySpan.appendChild(subcategoryTextNode);
    typeSpan.appendChild(document.createTextNode("Type: "));
    budgetSpan.appendChild(document.createTextNode("Budget: "));
    typeSpanValue.appendChild(typeTextNode);
    budgetSpanValue.appendChild(budgetTextNode);
    budgetSpanETB.appendChild(document.createTextNode(" ETB"));
    pDate.append(dateTextNode);
    

    div.className = "container-fluid project-card overflow-hidden";
    div.style = "cursor: pointer";
    imgDiv.className = "float-right overflow-hidden";
    image.className = "img-fluid rounded-circle";
    image.style = "width:35px; height:39px; padding-top: 4px;";
    if (container.ProfilePic != ""){
        image.src = "../assets/profile_pic/" + container.ProfilePic;
    }
    else{
        image.src = "../templates/images/user.png";
    }
    pTitleDiv.className = "overflow-hidden mt-sm-3";
    pTitle.className = "ml-sm-1 float-left project-title";
    pDescription.className = "ml-sm-1 mt-sm-3";
    pDescription.style = "font-weight: bold; color: #727272";
    pDescriptionValue.className = "ml-sm-3 mt-sm-2";
    typeDiv.className = "row overflow-hidden mt-sm-2";
    maskedDiv.className = "col-sm-5 masked ml-sm-1 overflow-hidden";
    unMaskedDiv.className = "col-sm-5 ml-sm-1 overflow-hidden";
    pCategory.className = "overflow-hidden";
    pSubcategory.className = "overflow-hidden";
    insideTypeDiv.className = "float-right";
    categorySpan.className = "float-left";
    subcategorySpan.className = "float-left";
    typeSpan.className = "color-change";
    budgetSpan.className = "color-change";
    pDate.style = "text-align: center; color: #da571f; font-weight: bold;";

    div.onclick = function(){
        window.location.href = "/View_Project?pid="+container.Project.ID;
    };
}


function clearCards(){
    var cards = document.getElementsByClassName("project-card");
    for (var i = 0; i < cards.length; i++) {
        cardElement = cards[i];
        cardElement.remove();
        console.log(i);
    }
}