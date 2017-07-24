
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
            setInterval(function(){ updateAdhocResultsTable(data); },3000);
        }
    });
}

function updateAdhocResultsTable(data) {
    $('#adhocResultTable').find('tr').remove().end();
    $('#adhocResultTable').append('<tr><th>Id</th><th>Agent</th><th>Time</th><th>Result</th><th>Graph</th></tr>');
    for(var idx in data.results) {
        $.ajax({
            url: "/api/results/"+data.results[idx]
        }).done(function(result){
            $('#adhocResultTable').append(
                $('<tr>').append(
                    $('<td>').append(result.ResultId)
                ).append(
                    $('<td>').append(result.AgentId)
                ).append(
                    $('<td>').append(result.AgentLabel)
                ).append(
                    $('<td>').append(result.AgentLocation)
                ).append(
                    $('<td>').append('<div id="graph_'+result.ResultId+'"></div>')
                )
            );
            var xAxis = [];
            var yAxis = [];

            for(var k in result.Datapoints) {
                xAxis.push(k);
                yAxis.push(result.Datapoints[k]);
            }

            var line = {
                x: xAxis,
                y: yAxis,
                type: 'scatter',
                dy: .1
            }

            var data = [line];

            console.log(xAxis);
            console.log(yAxis);

            Plotly.newPlot('graph',data);
        });
    }
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