package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Broadside</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center}
.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.rel{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.rel:hover{border-color:var(--leather)}
.rel-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.rel-title{font-size:.85rem;font-weight:700}
.rel-body{font-size:.7rem;color:var(--cd);margin-top:.2rem;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}
.rel-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}
.rel-actions{display:flex;gap:.3rem;flex-shrink:0}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.badge.draft{border-color:var(--gold);color:var(--gold)}.badge.published{border-color:var(--green);color:var(--green)}.badge.sent{border-color:var(--blue);color:var(--blue)}
.tag{font-size:.45rem;padding:.1rem .25rem;background:var(--bg3);color:var(--cd)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:500px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> BROADSIDE</h1><button class="btn btn-p" onclick="openForm()">+ New Release</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar"><input class="search" id="search" placeholder="Search releases..." oninput="render()"></div>
<div id="list"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/releases').then(function(r){return r.json()});items=r.releases||[];renderStats();render();}
function renderStats(){var total=items.length,draft=items.filter(function(r){return r.status==='draft'}).length,pub=items.filter(function(r){return r.status==='published'||r.status==='sent'}).length;
document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Total</div></div><div class="st"><div class="st-v">'+draft+'</div><div class="st-l">Drafts</div></div><div class="st"><div class="st-v" style="color:var(--green)">'+pub+'</div><div class="st-l">Published</div></div>';}
function render(){var q=(document.getElementById('search').value||'').toLowerCase();var f=items;
if(q)f=f.filter(function(r){return(r.title||'').toLowerCase().includes(q)||(r.outlet||'').toLowerCase().includes(q)||(r.contact||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No press releases.</div>';return;}
var h='';f.forEach(function(r){
h+='<div class="rel"><div class="rel-top"><div style="flex:1">';
h+='<div class="rel-title">'+esc(r.title)+'</div>';
if(r.body)h+='<div class="rel-body">'+esc(r.body)+'</div>';
h+='</div><div class="rel-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(''+r.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+r.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div><div class="rel-meta">';
if(r.status)h+='<span class="badge '+r.status+'">'+r.status+'</span>';
if(r.outlet)h+='<span>'+esc(r.outlet)+'</span>';
if(r.contact)h+='<span>'+esc(r.contact)+'</span>';
if(r.publish_date)h+='<span>'+r.publish_date+'</span>';
if(r.tags){r.tags.split(',').forEach(function(t){t=t.trim();if(t)h+='<span class="tag">#'+esc(t)+'</span>';});}
h+='</div></div>';});
document.getElementById('list').innerHTML=h;}
async function del(id){if(!confirm('Delete?'))return;await fetch(A+'/releases/'+id,{method:'DELETE'});load();}
function formHTML(rel){var i=rel||{title:'',body:'',contact:'',outlet:'',status:'draft',publish_date:'',tags:''};var isEdit=!!rel;
var h='<h2>'+(isEdit?'EDIT':'NEW')+' PRESS RELEASE</h2>';
h+='<div class="fr"><label>Title *</label><input id="f-title" value="'+esc(i.title)+'"></div>';
h+='<div class="fr"><label>Body</label><textarea id="f-body" rows="5">'+esc(i.body)+'</textarea></div>';
h+='<div class="row2"><div class="fr"><label>Contact</label><input id="f-contact" value="'+esc(i.contact)+'"></div>';
h+='<div class="fr"><label>Outlet</label><input id="f-outlet" value="'+esc(i.outlet)+'"></div></div>';
h+='<div class="row2"><div class="fr"><label>Status</label><select id="f-status">';
['draft','published','sent'].forEach(function(s){h+='<option value="'+s+'"'+(i.status===s?' selected':'')+'>'+s.charAt(0).toUpperCase()+s.slice(1)+'</option>';});
h+='</select></div><div class="fr"><label>Publish Date</label><input id="f-date" type="date" value="'+esc(i.publish_date)+'"></div></div>';
h+='<div class="fr"><label>Tags</label><input id="f-tags" value="'+esc(i.tags)+'" placeholder="comma separated"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Create')+'</button></div>';
return h;}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var r=null;for(var j=0;j<items.length;j++){if(items[j].id===id){r=items[j];break;}}if(!r)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(r);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var title=document.getElementById('f-title').value.trim();if(!title){alert('Title required');return;}
var body={title:title,body:document.getElementById('f-body').value.trim(),contact:document.getElementById('f-contact').value.trim(),outlet:document.getElementById('f-outlet').value.trim(),status:document.getElementById('f-status').value,publish_date:document.getElementById('f-date').value,tags:document.getElementById('f-tags').value.trim()};
if(editId){await fetch(A+'/releases/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/releases',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();}
function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
