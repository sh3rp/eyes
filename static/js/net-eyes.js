
function loadAgents() {
    $.ajax({
        url: "/api/agents"
    })
    .done(function(data){
        $('#viewPort').html(outputAgentTable(data));
        $("#allAgentCheckbox").click(function() {
            $(".agentCheckbox").prop("checked",$(".agentCheckbox").prop("checked"));
        });
    });
}

function outputAgentTable(data) {
    var htmlTable = '<table class="table table-striped">';
    htmlTable += '<thead><tr><th>' +
      '<input type="checkbox" aria-label="..." id="allAgentCheckbox" class="agentCheckbox">' + '</th><th>ID</th><th>Description</th><th>Location</th></tr></thead>';
    for(var key in data) {
        htmlTable += '<tr><td>' + 
        '<input type="checkbox" aria-label="..." class="agentCheckBox">';
        htmlTable += '</td><td>' + key + '</td><td>' + data[key].Label +'</td><td>' + data[key].Location + '</td></tr>';
    }
    htmlTable += '</table>';

    return htmlTable;
}

function addActiveSwitcher() {
    $('.navbar-nav li').click(function(e) {
        $('.navbar li.active').removeClass('active');
        var $this = $(this);
        if (!$this.hasClass('active')) {
            $this.addClass('active');
        }
        e.preventDefault();
        switch(e.target.id) {
            case "agents": // load agents
                loadAgents();
                break;
            case "schedules": // load schedules
                alert("Loading schedules.");
                break;
            case "visuals": // load visuals
                alert("Loading visuals.");
                break;
        }
    });
}

$(document).ready(function() {
    addActiveSwitcher();
});