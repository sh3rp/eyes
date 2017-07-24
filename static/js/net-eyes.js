
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
    for(var idx in data.results) {
        $.ajax({
            url: "/api/results/"+data.results[idx]
        }).done(function(result){
            $('#adhocGraphTitle').text(result.TargetHost);
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
                dy: .1,
                y0: yAxis[0] - .3
            }

        var layout = {
            autosize: false,
            width: 500,
            height: 500,
            margin: {
                l: 50,
                r: 50,
                b: 100,
                t: 100,
                pad: 4
            },
            yaxis: {
                autotick: false,
                ticks: 'outside',
                tick0: 0,
                dtick: 0.25,
                ticklen: 2,
                tickwidth: 1,
                tickcolor: '#000'
            },
            paper_bgcolor: '#ffffff',
            plot_bgcolor: '#ffffff'
        };

            var data = [line];

            console.log(xAxis);
            console.log(yAxis);

            Plotly.newPlot('graph',data,layout);
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