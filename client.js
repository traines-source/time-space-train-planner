currentSelected = undefined;
function selectEdge(edgeId) {
    setColorForDependents(currentSelected, 'black');
    setColorForDependents(edgeId, 'red');
    currentSelected = edgeId;
}

function setColorForDependents(edgeId, color) {
    if (edgeId == undefined) {
        return;
    }
    document.getElementById(edgeId).style.stroke = color;
    const dependents = document.getElementsByClassName('shortest-path-'+edgeId);
    for(let i=0;i<dependents.length;i++) {
        dependents[i].style.stroke = color;
    }
}


const edges = document.getElementsByClassName('edge');
for(let i=0;i<edges.length;i++) {
    edges[i].onclick = function (evt) {
        console.log('selected ', this.id);
        selectEdge(this.id);
    }
}
