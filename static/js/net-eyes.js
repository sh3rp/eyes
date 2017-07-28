
function loadAdhocScreen() {
    $('#adhocForm').attr('hidden',false);
    $('#adhocGraph').attr('hidden',false);
    $.ajax({
        url:"/api/agents"
    }).done(function(data) {
        $('#adhocAgentList').find('option').remove().end();
        for(var key in data) {
            $('#adhocAgentList').append($('<option>', {
                value: key,
                text: data[key].Info.label + " (" + data[key].Info.location + ")"
            }));
        }
    })
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

function postAdhocRequest() {
    $.ajax({
        type: "POST",
        url: "/api/agent.control",
        data: JSON.stringify({
            Host: $('#adhocHost').val(),
            Type: $('#adhocType').val(),
            Agents: $('#adhocAgentList').val(),
            MaxPoints: parseInt($('#adhocMaxPoints').val())
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