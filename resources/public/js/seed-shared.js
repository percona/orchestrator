
function appendSeedDetails(seed, selector) {    	
	var row = '<tr>';
	var statusMessage;
	if (seed.IsComplete) {
		statusMessage = (seed.IsSuccessful ? '<span class="text-success">Success</span>' : '<span class="text-danger">Fail</span>');
	} else {
		statusMessage = '<span class="text-info">Active</span>';
	}
	row += '<td>' + statusMessage + '</td>';
	row += '<td><a href="' + appUrl('/web/seed-details/' + encodeURIComponent(seed.SeedId)) + '">' + escapeHtml(seed.SeedId) + '</a></td>';
	row += '<td><a href="' + appUrl('/web/agent/'+encodeURIComponent(seed.TargetHostname)) + '">'+escapeHtml(seed.TargetHostname)+'</a></td>';
	row += '<td><a href="' + appUrl('/web/agent/'+encodeURIComponent(seed.SourceHostname)) + '">'+escapeHtml(seed.SourceHostname)+'</a></td>';
	row += '<td>' + escapeHtml(seed.StartTimestamp) + '</td>';
	row += '<td>' + (seed.IsComplete ? escapeHtml(seed.EndTimestamp) : '<button class="btn btn-xs btn-danger" data-command="abort-seed" data-seed-source-host="'+escapeHtml(seed.SourceHostname)+'" data-seed-target-host="'+escapeHtml(seed.TargetHostname)+'" data-seed-id="' + escapeHtml(seed.SeedId) + '">Abort</button>') + '</td>';
	row += '</tr>';
	$(selector).append(row);
    hideLoader();
}

function appendSeedState(seedState) {    	
	var action = seedState.Action;
	action = action.replace(/Copied ([0-9]+).([0-9]+) bytes (.*$)/, function(match, match1, match2, match3) { 
		return "Copied " + toHumanFormat(match1) + " / " + toHumanFormat(match2) + " " + match3;
	});
	var row = '<tr>';
	row += '<td>' + escapeHtml(seedState.StateTimestamp) + '</td>';
	row += '<td>' + escapeHtml(action) + '</td>';
	row += '<td>' + escapeHtml(seedState.ErrorMessage) + '</td>';
	row += '</tr>';
	$("[data-agent=seed_states]").append(row);
    hideLoader();
}

$("body").on("click", "button[data-command=abort-seed]", function(event) {
	var seedId = $(event.target).attr("data-seed-id");
	var sourceHost = $(event.target).attr("data-seed-source-host");
	var targetHost = $(event.target).attr("data-seed-target-host");

	var message = "Are you sure you wish to abort seed " + escapeHtml(seedId) + " from <code><strong>" + 
		escapeHtml(sourceHost) + "</strong></code> to <code><strong>" + 
		escapeHtml(targetHost) + "</strong></code> ?";
	bootbox.confirm(message, function(confirm) {
		if (confirm) {
	    	showLoader();
	        $.get(appUrl("/api/agent-abort-seed/"+encodeURIComponent(seedId)), function (operationResult) {
				hideLoader();
				if (operationResult.Code == "ERROR") {
					addAlert(escapeHtml(operationResult.Message))
				} else {
					location.reload();
				}	
	        }, "json").fail(function (operationResult) {
				hideLoader();
				addAlert(escapeHtml((operationResult.responseJSON && operationResult.responseJSON.Message) || "Request failed"))
			  });
		}
	});
});
