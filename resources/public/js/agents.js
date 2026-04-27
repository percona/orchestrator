
$(document).ready(function () {
    showLoader();
    activateRefreshTimer();
    
    $.get(appUrl("/api/agents"), function (agents) {
    	displayAgents(agents);
    }, "json");
    function displayAgents(agents) {
        hideLoader();
        
        agents.forEach(function (agent) {
    		$("#agents").append('<div xmlns="http://www.w3.org/1999/xhtml" class="popover instance right" data-agent-name="'+escapeHtml(agent.Hostname)+'"><div class="arrow"></div><div class="popover-content"></div></div>');
    		var popoverElement = $("#agents .popover.instance.right").last();
    		//var title = agent.Hostname;
    		//popoverElement.find("h3 a").html(title);
    	    var contentHtml = ''
    	    	+ '<a href="' + appUrl('/web/agent/'+ encodeURIComponent(agent.Hostname)) +'" class="small">'
    	    	+ escapeHtml(agent.Hostname)
    	    	+ '</a>'
    			;
    	    popoverElement.find(".popover-content").html(contentHtml);
        });     
        
        $("div.popover").popover();
        $("div.popover").show();
	
        if (agents.length == 0) {
        	addAlert("No agents found");
        }
    }
});	
