const xmlns = "http://www.w3.org/2000/svg";

currentSelected = undefined;
function selectEdge(edgeId) {
    setSelectedForDependents(currentSelected, false);
    setSelectedForDependents(edgeId, true);
    currentSelected = edgeId;
}

function setActive(edgeId) {
    document.getElementById('details').innerHTML = document.getElementById(edgeId+'-label').children[1].innerHTML;
    document.getElementById('legend').style.display = 'none';
}

function setInactive() {
    document.getElementById('details').innerHTML = '';
    document.getElementById('legend').style.display = 'block';
}

function setSelectedForDependents(edgeId, selected) {
    if (edgeId == undefined) {
        return;
    }
    const dependents = document.getElementsByClassName('sp-'+edgeId);
    for(let i=0;i<dependents.length;i++) {
        setSelectedForElementWithLabels(dependents[i], selected);
        setPreviousNextForElement(dependents[i], selected);
    }
}

function setSelectedForElementWithLabels(element, selected) {
    setSelectedForElement(element, selected);
    const labels = document.getElementsByClassName('label-'+element.id);
    for(let i=0;i<labels.length;i++) {
        setSelectedForElement(labels[i], selected);
    }
}

function setSelectedForElement(element, selected) {
    if (selected) {
        element.className.baseVal += " selected";
    } else {
        element.className.baseVal =  element.className.baseVal.replace(" selected", "");
    }    
}

function setPreviousNextForElement(element, selected) {
    if (selected) {
        if (!element.getAttribute('d')) {
            return;
        }
        const coords = element.getAttribute('d').split(/[^0-9]+/);
        const from = {x: parseFloat(coords[1]), y: parseFloat(coords[2])};
        const to = {x: parseFloat(coords[3]), y: parseFloat(coords[4])};
        const margin = 25;
        createArrow(element, true, element.dataset.pa, to.x, to.y-margin);
        createArrow(element, true, element.dataset.pd, from.x, from.y-margin);
        createArrow(element, false, element.dataset.na, to.x, to.y+margin);
        createArrow(element, false, element.dataset.nd, from.x, from.y+margin);
    } else {        
        const arrows = document.getElementsByClassName('previous-next-arrow');
        for(let i=arrows.length-1;i>=0;i--) {
            arrows[i].remove();
        }
    }
}

function createArrow(parent, previous, targetId, x, y) {
    if (targetId == undefined || targetId == '') {
        return;
    }
    const e = document.createElementNS(xmlns, 'text');
    e.className.baseVal = 'previous-next-arrow';
    e.setAttribute('x', x);
    e.setAttribute('y', y);
    e.innerHTML = previous ? '▲' : '▼';
    e.onclick = function (evt) {
        console.log('arrow select ', targetId);
        selectEdge(targetId);
        //setActive(targetId);
    };
    parent.parentNode.appendChild(e);
}

function selectListener(evt) {
    const id = this.id.replace('-toucharea', '');
    console.log('selected ', id);
    selectEdge(id);
    //setActive(id);
}

const edges = document.getElementsByClassName('edge-toucharea');
for(let i=0;i<edges.length;i++) {
    edges[i].onclick = selectListener; 
}