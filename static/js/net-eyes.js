
function loadAdhocScreen() {
    $('#adhocForm').attr('hidden',false);
    $('#adhocResults').attr('hidden',false);
    $.ajax({
        url:"/api/agents"
    }).done(function(data) {
        $('#adhocAgentList').find('option').remove().end();
        for(var key in data) {
            $('#adhocAgentList').append($('<option>', {
                value: key,
                text: data[key].Label
            }));
        }
    })
}

function postAdhocRequest() {
    $.ajax({
        type: "POST",
        url: "/api/agent.control",
        data: JSON.stringify({
            Host: $('#adhocHost').val(),
            Type: $('#adhocType').val(),
            Agents: $('#adhocAgentList').val()
        }),
        dataType: "json"
    }).done(function(data) {
        if(data.code == 0) {
            for(var id in data.results) {
                console.log("Result: " + data.results[id]);
            }
        }
    });
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
            case "adhocProbe":
                loadAdhocScreen();
            case "agents": // load agents
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
    $('#adhocSubmit').click(function(e) {
        postAdhocRequest();
    });
});