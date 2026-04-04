package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Podium</title>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital,wght@0,400;0,700;1,400&family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:800px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center;font-family:var(--mono);cursor:pointer;transition:border-color .2s}
.st:hover,.st.active{border-color:var(--rust)}.st.active .st-v{color:var(--rust)}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center;flex-wrap:wrap}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.fb{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem;margin-bottom:.5rem;display:flex;gap:.8rem;transition:border-color .2s}
.fb:hover{border-color:var(--leather)}
.vote-col{display:flex;flex-direction:column;align-items:center;min-width:40px;flex-shrink:0}
.vote-btn{background:none;border:none;color:var(--cm);cursor:pointer;font-size:1.1rem;padding:.1rem;line-height:1;transition:color .15s}
.vote-btn:hover{color:var(--rust)}
.vote-count{font-family:var(--mono);font-size:.95rem;font-weight:700;color:var(--cream)}
.fb-content{flex:1;min-width:0}
.fb-title{font-size:.88rem;margin-bottom:.1rem}
.fb-body{font-size:.72rem;color:var(--cd);margin-top:.2rem}
.fb-meta{font-family:var(--mono);font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}
.fb-actions{display:flex;gap:.3rem;flex-shrink:0;align-self:flex-start}
.badge{font-family:var(--mono);font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.badge.open{border-color:var(--blue);color:var(--blue)}.badge.planned{border-color:var(--gold);color:var(--gold)}.badge.in-progress{border-color:var(--rust);color:var(--rust)}.badge.done{border-color:var(--green);color:var(--green)}.badge.closed{border-color:var(--cm);color:var(--cm)}
.cat-badge{font-family:var(--mono);font-size:.5rem;padding:.1rem .3rem;background:var(--bg3);color:var(--cd)}
.tag{font-family:var(--mono);font-size:.5rem;padding:.1rem .3rem;background:var(--bg3);color:var(--cm)}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}.row2{grid-template-columns:1fr}.toolbar{flex-direction:column}.search{min-width:100%}}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> PODIUM</h1><button class="btn btn-p" onclick="openForm()">+ New Feedback</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search feedback..." oninput="render()">
<select class="filter-sel" id="status-filter" onchange="render()"><option value="">All Status</option><option value="open">Open</option><option value="planned">Planned</option><option value="in-progress">In Progress</option><option value="done">Done</option><option value="closed">Closed</option></select>
<select class="filter-sel" id="sort-sel" onchange="render()"><option value="votes">Top Voted</option><option value="newest">Newest</option><option value="oldest">Oldest</option></select>
</div>
<div id="list"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api/feedback',items=[],filter='all',editId=null;

async function load(){var r=await fetch(A).then(function(r){return r.json()});items=r.feedback||[];renderStats();render();}

function renderStats(){
var total=items.length;
var open=items.filter(function(i){return i.status==='open'}).length;
var planned=items.filter(function(i){return i.status==='planned'}).length;
var done=items.filter(function(i){return i.status==='done'}).length;
document.getElementById('stats').innerHTML=[
{l:'Total',v:total,f:'all'},{l:'Open',v:open,f:'open'},{l:'Planned',v:planned,f:'planned'},{l:'Done',v:done,f:'done'}
].map(function(x){return '<div class="st'+(filter===x.f?' active':'')+'" onclick="setFilter(''+x.f+'')"><div class="st-v">'+x.v+'</div><div class="st-l">'+x.l+'</div></div>'}).join('');
}

function setFilter(f){filter=f;renderStats();render();}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var sf=document.getElementById('status-filter').value;
var sort=document.getElementById('sort-sel').value;
var f=items.slice();
if(filter!=='all')f=f.filter(function(i){return i.status===filter});
if(sf)f=f.filter(function(i){return i.status===sf});
if(q)f=f.filter(function(i){return(i.title||'').toLowerCase().includes(q)||(i.body||'').toLowerCase().includes(q)||(i.author||'').toLowerCase().includes(q)});
if(sort==='votes')f.sort(function(a,b){return(b.votes||0)-(a.votes||0)});
else if(sort==='newest')f.sort(function(a,b){return(b.created_at||'').localeCompare(a.created_at||'')});
else f.sort(function(a,b){return(a.created_at||'').localeCompare(b.created_at||'')});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No feedback yet. Be the first to suggest something.</div>';return;}
var h='';f.forEach(function(i){
var statusCls=(i.status||'open').replace(/_/g,'-');
h+='<div class="fb"><div class="vote-col">';
h+='<button class="vote-btn" onclick="upvote(''+i.id+'')">&#9650;</button>';
h+='<div class="vote-count">'+(i.votes||0)+'</div>';
h+='</div><div class="fb-content">';
h+='<div class="fb-title">'+esc(i.title)+'</div>';
if(i.body)h+='<div class="fb-body">'+esc(i.body)+'</div>';
h+='<div class="fb-meta">';
h+='<span class="badge '+statusCls+'">'+esc(i.status||'open')+'</span>';
if(i.category)h+='<span class="cat-badge">'+esc(i.category)+'</span>';
if(i.author)h+='<span>'+esc(i.author)+'</span>';
if(i.tags){i.tags.split(',').forEach(function(t){t=t.trim();if(t)h+='<span class="tag">#'+esc(t)+'</span>';});}
h+='<span>'+ft(i.created_at)+'</span>';
h+='</div></div>';
h+='<div class="fb-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(''+i.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+i.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
});
document.getElementById('list').innerHTML=h;
}

async function upvote(id){
var item=null;for(var j=0;j<items.length;j++){if(items[j].id===id){item=items[j];break;}}
if(!item)return;
await fetch(A+'/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({votes:(item.votes||0)+1})});
load();
}

async function del(id){if(!confirm('Delete this feedback?'))return;await fetch(A+'/'+id,{method:'DELETE'});load();}

function formHTML(item){
var i=item||{title:'',body:'',author:'',category:'',status:'open',tags:''};
var isEdit=!!item;
var h='<h2>'+(isEdit?'EDIT FEEDBACK':'NEW FEEDBACK')+'</h2>';
h+='<div class="fr"><label>Title *</label><input id="f-title" value="'+esc(i.title)+'" placeholder="Short summary of your feedback"></div>';
h+='<div class="fr"><label>Description</label><textarea id="f-body" rows="4" placeholder="Describe your idea or issue in detail...">'+esc(i.body)+'</textarea></div>';
h+='<div class="row2"><div class="fr"><label>Author</label><input id="f-author" value="'+esc(i.author)+'" placeholder="Your name"></div>';
h+='<div class="fr"><label>Category</label><input id="f-cat" value="'+esc(i.category)+'" placeholder="e.g. feature, bug, ux"></div></div>';
h+='<div class="row2"><div class="fr"><label>Status</label><select id="f-status">';
['open','planned','in-progress','done','closed'].forEach(function(s){h+='<option value="'+s+'"'+(i.status===s?' selected':'')+'>'+s.charAt(0).toUpperCase()+s.slice(1)+'</option>';});
h+='</select></div><div class="fr"><label>Tags</label><input id="f-tags" value="'+esc(i.tags)+'" placeholder="comma separated"></div></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Submit')+'</button></div>';
return h;
}

function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');document.getElementById('f-title').focus();}
function openEdit(id){var item=null;for(var j=0;j<items.length;j++){if(items[j].id===id){item=items[j];break;}}if(!item)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(item);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}

async function submit(){
var title=document.getElementById('f-title').value.trim();
if(!title){alert('Title is required');return;}
var body={title:title,body:document.getElementById('f-body').value.trim(),author:document.getElementById('f-author').value.trim(),category:document.getElementById('f-cat').value.trim(),status:document.getElementById('f-status').value,tags:document.getElementById('f-tags').value.trim()};
if(editId){await fetch(A+'/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{body.votes=0;await fetch(A,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();
}

function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
