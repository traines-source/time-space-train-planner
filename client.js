currentSelected = undefined;
function selectEdge(edgeId) {
    setSelectedForDependents(currentSelected, false);
    setSelectedForDependents(edgeId, true);
    currentSelected = edgeId;
}

function setSelectedForDependents(edgeId, selected) {
    if (edgeId == undefined) {
        return;
    }
    setSelectedForElement(document.getElementById(edgeId), selected);
    const dependents = document.getElementsByClassName('shortest-path-'+edgeId);
    for(let i=0;i<dependents.length;i++) {
        setSelectedForElement(dependents[i], selected);
    }
}

function setSelectedForElement(element, selected) {
    if (selected) {
        element.className.baseVal += " selected";
    } else {
        element.className.baseVal = element.className.baseVal.replace(" selected", "");
    }
}


const edges = document.getElementsByClassName('edge-toucharea');
for(let i=0;i<edges.length;i++) {
    edges[i].onclick = function (evt) {
        const id = this.id.replace('-toucharea', '');
        console.log('selected ', id);
        selectEdge(id);
    }
}
