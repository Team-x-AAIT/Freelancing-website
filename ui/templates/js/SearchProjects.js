window.addEventListener('load', onPageLoad);

function onPageLoad() {

    var csrfError = document.getElementById("csrf_error");
    var csrfErrorDiv = document.getElementById("csrf_error_div");
    if (csrfError.value != "") {
        csrfErrorDiv.style.display = "block";
    }
}

function onViewProject(project){

    var pid = project.id;
    window.location.href = "/View_Project?pid=" + pid;

}