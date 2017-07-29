
function loadAgents() {
    $.ajax({
        url: '/api/agents'
    }).done(function(data){
        $('#agentList').find('option').remove().end();
        for(var idx in data) {
            $('#agentList').append($('<option>', {
                value: idx,
                text: data[idx].Info.label + " (" + data[idx].Info.location + ")"
            }));
        }
    });
}

function loadAdhocScreen() {
    $('#viewPort').load('/html/adhocProbe.html',function(){
        $('#adhocType').change(function(){
        var value = $(this).val();
        switch(value) {
            case "TCP":
               $('#adhocTCPTypeOptions').attr('hidden',false);
               break;
            default:
                $('#adhocTCPTypeOptions').attr('hidden',true);
                break;
        }
        });
    });
    loadAgents();
}

function loadAgentScreen() {
    $('#viewPort').load('/html/agents.html');
    $.ajax({
        url:"/api/agents"
    }).done(function(data) {
        $('#agents table').append('<tr><th>Name</th><th>Location</th><th>OS</th><th>Version</th></tr>');
        // add rows to table here
        for(var key in data) {
            var d = data[key];
            $('#agents table').append('<tr><td>'+d.Info.label+'</td>'+ 
                '<td>'+d.Info.location+'</td>'+
                '<td>'+d.Info.os+'</td>'+
                '<td>'+d.Info.agentVersion+'</td>'+
                '</tr>');
        }
    });
}

function postCancelRequest(id) {
    $.ajax({
        type: "POST",
        url: "/api/agent.cancel/"+id,
    }).done(function(data) {
        $('#adhocGraph').attr('hidden',true);
        $("[id^='adhocGraphImage_']").remove();
        clearInterval(intervalId);
    });
}

var intervalId;

function getOptions() {
    var options = {};

    if($('#adhocTCPPort').length) {
        options['port'] = $('#adhocTCPPort').val();
    }

    console.log('options = ' + options);

    return options;
}

function postAdhocRequest() {
    $.ajax({
        type: "POST",
        url: "/api/agent.control",
        data: JSON.stringify({
            Host: $('#adhocHost').val(),
            Type: $('#adhocType').val(),
            Agents: $('#agentList').val(),
            MaxPoints: parseInt($('#adhocMaxPoints').val()),
            Options: getOptions()
        }),
        dataType: "json"
    }).done(function(data) {
        $('#adhocCancel').click(function() {
            for(var id in data.results) {
                postCancelRequest(data.results[id]);
            }
        });
        if(data.code == 0) {
            $('#adhocGraph').attr('hidden',false);
            for(var idx in data.results) {
                $('#adhocGraphImage').append('<div id="adhocGraphImage_' + data.results[idx]+'"></div>');
            }
            intervalId = setInterval(function(){ updateAdhocResultsTable(data); },1000);
        }
    });
}

function updateAdhocResultsTable(data) {
    for(var idx in data.results) {
        $.ajax({
            url: "/api/results/"+data.results[idx]
        }).done(function(result){
            $('#adhocGraphTitle').text(result.TargetHost);
            var data = [];

            for(var k in result.Datapoints) {
                var point = {'time':k,'latency':result.Datapoints[k]};
                data.push(point);
            }  

            MG.data_graphic({
                title: 'Latency ('+result.AgentLocation+')',
                description: 'Latency to ' + result.TargetHost + ' from ' + result.AgentLocation + '.',
                data: data, // an array of objects, such as [{value:100,date:...},...]
                width: 600,
                height: 250,
                target: '#adhocGraphImage_'+result.ResultId, // the html element that the graphic is inserted in
                x_accessor: 'time',  // the key that accesses the x value
                y_accessor: 'latency', // the key that accesses the y value
                transition_on_update: true
            });

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
                break;
            case "agents": // load agents
                loadAgentScreen();
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

var loadAgentsIntervalId = null;

$(document).ready(function() {
    addActiveSwitcher();
});