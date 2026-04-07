package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Podium</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--purple:#9d6bb8;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;font-size:13px}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.2rem 1.5rem;max-width:1000px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700;color:var(--gold)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.list{display:flex;flex-direction:column;gap:.5rem}
.item{background:var(--bg2);border:1px solid var(--bg3);padding:.9rem 1rem;display:flex;gap:.9rem;align-items:flex-start;transition:border-color .15s}
.item:hover{border-color:var(--leather)}
.item.completed,.item.closed{opacity:.55}
.vote-col{display:flex;flex-direction:column;align-items:center;gap:.2rem;flex-shrink:0;min-width:54px}
.vote-btn{width:100%;padding:.3rem .4rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cd);cursor:pointer;font-family:var(--mono);font-size:.65rem;transition:.15s;display:flex;flex-direction:column;align-items:center;gap:.1rem}
.vote-btn:hover{border-color:var(--rust);color:var(--rust)}
.vote-btn.voting{opacity:.5;pointer-events:none}
.vote-arrow{font-size:.85rem;line-height:1}
.vote-count{font-size:1rem;font-weight:700;color:var(--gold);font-family:var(--mono)}
.vote-down{font-size:.5rem;color:var(--cm);background:none;border:none;cursor:pointer;padding:.1rem .3rem}
.vote-down:hover{color:var(--red)}
.body-col{flex:1;min-width:0;cursor:pointer}
.item-title{font-size:.85rem;font-weight:700;color:var(--cream)}
.item-body{font-size:.7rem;color:var(--cd);margin-top:.3rem;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}
.item-meta{font-size:.55rem;color:var(--cm);margin-top:.4rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm);font-weight:700}
.badge.open{border-color:var(--blue);color:var(--blue)}
.badge.planned{border-color:var(--purple);color:var(--purple)}
.badge.in_progress{border-color:var(--orange);color:var(--orange)}
.badge.completed{border-color:var(--green);color:var(--green)}
.badge.closed{border-color:var(--cm);color:var(--cm)}
.badge.cat{border-color:var(--leather);color:var(--leather)}
.tag-list{display:flex;gap:.3rem;flex-wrap:wrap}
.tag{font-size:.5rem;padding:.05rem .3rem;background:var(--bg3);color:var(--cd);font-family:var(--mono)}
.item-extra{font-size:.55rem;color:var(--cd);margin-top:.4rem;padding-top:.3rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.item-extra-row{display:flex;gap:.4rem}
.item-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.item-extra-val{color:var(--cream)}
.btn{font-family:var(--mono);font-size:.6rem;padding:.3rem .55rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:520px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto;color:var(--red);border-color:#3a1a1a}
.acts .btn-del:hover{border-color:var(--red);color:var(--red)}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> PODIUM</h1>
<button class="btn btn-p" onclick="openNew()">+ Submit Idea</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search title, body, tags..." oninput="debouncedRender()">
<select class="filter-sel" id="status-filter" onchange="render()">
<option value="">All Statuses</option>
<option value="open">Open</option>
<option value="planned">Planned</option>
<option value="in_progress">In Progress</option>
<option value="completed">Completed</option>
<option value="closed">Closed</option>
</select>
<select class="filter-sel" id="category-filter" onchange="render()">
<option value="">All Categories</option>
</select>
<select class="filter-sel" id="sort-filter" onchange="render()">
<option value="votes">Most Voted</option>
<option value="newest">Newest</option>
<option value="oldest">Oldest</option>
</select>
</div>
<div id="list" class="list"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='feedback';

var fields=[
{name:'title',label:'Title',type:'text',required:true},
{name:'body',label:'Description',type:'textarea'},
{name:'author',label:'Author',type:'text'},
{name:'category',label:'Category',type:'select_or_text',options:[]},
{name:'status',label:'Status',type:'select',options:['open','planned','in_progress','completed','closed']},
{name:'tags',label:'Tags',type:'text',placeholder:'comma separated'}
];

var items=[],itemExtras={},editId=null,searchTimer=null;

function fmtDate(s){
if(!s)return'';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return s;
return d.toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'});
}catch(e){return s}
}

function fieldByName(n){
for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];
return null;
}

function debouncedRender(){
clearTimeout(searchTimer);
searchTimer=setTimeout(render,200);
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var resps=await Promise.all([
fetch(A+'/feedback').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
items=resps[0].feedback||[];
renderStats(resps[1]||{});

try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
itemExtras=ex||{};
items.forEach(function(it){
var x=itemExtras[it.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(it[k]===undefined)it[k]=x[k]});
});
}catch(e){itemExtras={}}

populateCategoryFilter();
}catch(e){
console.error('load failed',e);
items=[];
}
render();
}

function populateCategoryFilter(){
var sel=document.getElementById('category-filter');
if(!sel)return;
var current=sel.value;
var seen={};
var cats=[];
items.forEach(function(it){
if(it.category&&!seen[it.category]){seen[it.category]=true;cats.push(it.category)}
});
cats.sort();
sel.innerHTML='<option value="">All Categories</option>'+cats.map(function(c){return'<option value="'+esc(c)+'"'+(c===current?' selected':'')+'>'+esc(c)+'</option>'}).join('');
}

function renderStats(s){
var total=s.total||0;
var totalVotes=s.total_votes||0;
var topVotes=s.top_votes||0;
var byStatus=s.by_status||{};
var open=byStatus.open||0;
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Ideas</div></div>'+
'<div class="st"><div class="st-v">'+open+'</div><div class="st-l">Open</div></div>'+
'<div class="st"><div class="st-v">'+totalVotes+'</div><div class="st-l">Total Votes</div></div>'+
'<div class="st"><div class="st-v">'+topVotes+'</div><div class="st-l">Top Vote</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var sf=document.getElementById('status-filter').value;
var cf=document.getElementById('category-filter').value;
var sort=document.getElementById('sort-filter').value;

var f=items.slice();
if(q)f=f.filter(function(it){
return(it.title||'').toLowerCase().includes(q)||
(it.body||'').toLowerCase().includes(q)||
(it.tags||'').toLowerCase().includes(q);
});
if(sf)f=f.filter(function(it){return it.status===sf});
if(cf)f=f.filter(function(it){return it.category===cf});

if(sort==='newest'){
f.sort(function(a,b){return(b.created_at||'').localeCompare(a.created_at||'')});
}else if(sort==='oldest'){
f.sort(function(a,b){return(a.created_at||'').localeCompare(b.created_at||'')});
}else{
f.sort(function(a,b){
if((b.votes||0)!==(a.votes||0))return(b.votes||0)-(a.votes||0);
return(b.created_at||'').localeCompare(a.created_at||'');
});
}

if(!f.length){
var msg=window._emptyMsg||'No ideas yet. Be the first to submit one.';
document.getElementById('list').innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}

var h='';
f.forEach(function(it){h+=itemHTML(it)});
document.getElementById('list').innerHTML=h;
}

function itemHTML(it){
var cls='item '+(it.status||'open');

var h='<div class="'+cls+'">';

// Vote column
h+='<div class="vote-col">';
h+='<button class="vote-btn" data-id="'+esc(it.id)+'" onclick="upvote(\''+esc(it.id)+'\',event)">';
h+='<span class="vote-arrow">&#9650;</span>';
h+='<span class="vote-count">'+(it.votes||0)+'</span>';
h+='</button>';
h+='<button class="vote-down" onclick="downvote(\''+esc(it.id)+'\',event)" title="Downvote">&#9660;</button>';
h+='</div>';

// Body column
h+='<div class="body-col" onclick="openEdit(\''+esc(it.id)+'\')">';
h+='<div class="item-title">'+esc(it.title)+'</div>';
if(it.body)h+='<div class="item-body">'+esc(it.body)+'</div>';

h+='<div class="item-meta">';
if(it.status)h+='<span class="badge '+esc(it.status)+'">'+esc(it.status.replace(/_/g,' '))+'</span>';
if(it.category)h+='<span class="badge cat">'+esc(it.category)+'</span>';
if(it.author)h+='<span>by '+esc(it.author)+'</span>';
if(it.created_at)h+='<span>'+esc(fmtDate(it.created_at))+'</span>';
h+='</div>';

if(it.tags){
var tagList=String(it.tags).split(',').map(function(t){return t.trim()}).filter(function(t){return t});
if(tagList.length){
h+='<div class="tag-list" style="margin-top:.4rem">';
tagList.forEach(function(t){h+='<span class="tag">#'+esc(t)+'</span>'});
h+='</div>';
}
}

// Custom field display
var customRows='';
fields.forEach(function(f){
if(!f.isCustom)return;
var v=it[f.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="item-extra-row">';
customRows+='<span class="item-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="item-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="item-extra">'+customRows+'</div>';

h+='</div></div>';
return h;
}

// ─── Voting ──────────────────────────────────────────────────────

async function upvote(id,ev){
ev.stopPropagation();
var btn=ev.currentTarget;
btn.classList.add('voting');
try{
var r=await fetch(A+'/feedback/'+id+'/upvote',{method:'POST'});
if(r.ok){
var d=await r.json();
// Update in-memory and re-render
for(var i=0;i<items.length;i++)if(items[i].id===id){items[i].votes=d.votes;break}
load(); // Re-fetch stats
}
}catch(e){}
btn.classList.remove('voting');
}

async function downvote(id,ev){
ev.stopPropagation();
try{
var r=await fetch(A+'/feedback/'+id+'/downvote',{method:'POST'});
if(r.ok){
var d=await r.json();
for(var i=0;i<items.length;i++)if(items[i].id===id){items[i].votes=d.votes;break}
load();
}
}catch(e){}
}

// ─── Modal ────────────────────────────────────────────────────────

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph=f.placeholder?(' placeholder="'+esc(f.placeholder)+'"'):'';
var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
var disp=String(o).charAt(0).toUpperCase()+String(o).slice(1).replace(/_/g,' ');
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(disp)+'</option>';
});
h+='</select>';
}else if(f.type==='select_or_text'){
h+='<input list="dl-'+f.name+'" type="text" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
h+='<datalist id="dl-'+f.name+'">';
var opts=(f.options||[]).slice();
items.forEach(function(itm){
if(itm.category&&opts.indexOf(itm.category)===-1)opts.push(itm.category);
});
opts.forEach(function(o){h+='<option value="'+esc(String(o))+'">'});
h+='</datalist>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="3"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='number'){
h+='<input type="number" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}else{
var inputType=f.type||'text';
h+='<input type="'+esc(inputType)+'" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}
h+='</div>';
return h;
}

function formHTML(item){
var it=item||{};
var isEdit=!!item;
var h='<h2>'+(isEdit?'EDIT IDEA':'NEW IDEA')+'</h2>';

h+=fieldHTML(fieldByName('title'),it.title);
h+=fieldHTML(fieldByName('body'),it.body);
h+='<div class="row2">'+fieldHTML(fieldByName('author'),it.author)+fieldHTML(fieldByName('category'),it.category)+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('status'),it.status||'open')+fieldHTML(fieldByName('tags'),it.tags)+'</div>';

var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var label=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,it[f.name])});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit){
h+='<button class="btn btn-del" onclick="delItem()">Delete</button>';
}
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Submit')+'</button>';
h+='</div>';
return h;
}

function openNew(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var t=document.getElementById('f-title');
if(t)t.focus();
}

function openEdit(id){
var it=null;
for(var i=0;i<items.length;i++){if(items[i].id===id){it=items[i];break}}
if(!it)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(it);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

async function submit(){
var titleEl=document.getElementById('f-title');
if(!titleEl||!titleEl.value.trim()){alert('Title is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val;
if(f.type==='number')val=parseFloat(el.value)||0;
else val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/feedback/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/feedback',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Submit failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}
closeModal();
load();
}

async function delItem(){
if(!editId)return;
if(!confirm('Delete this idea?'))return;
await fetch(A+'/feedback/'+editId,{method:'DELETE'});
closeModal();
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.categories)){
var catField=fieldByName('category');
if(catField)catField.options=cfg.categories;
}

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
if(fieldByName(cf.name))return;
fields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[],
isCustom:true
});
});
}
}).catch(function(){
}).finally(function(){
load();
});
})();
</script>
</body>
</html>`
