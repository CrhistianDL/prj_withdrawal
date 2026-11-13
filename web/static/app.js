let operacion='ABONO';
const $=id=>document.getElementById(id);
const dinero=n=>new Intl.NumberFormat('es-PE',{style:'currency',currency:'PEN'}).format(n);
const safe=t=>String(t).replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
async function obtener(url,opciones){const r=await fetch(url,opciones);const d=await r.json();if(!r.ok)throw Error(d.error||'Error de conexión');return d}
async function actualizar(){try{const e=await obtener('/api/estado');$('titular').textContent=e.titular;$('saldo').textContent=dinero(e.saldo);const datos=await obtener('/api/transacciones');$('historial').innerHTML=datos.length?datos.slice().reverse().map(t=>`<tr><td>${safe(new Date(t.fecha).toLocaleString('es-PE'))}</td><td>${safe(t.operacion)}</td><td>${dinero(t.monto)}</td><td><span class="pill ${safe(t.estado)}">${safe(t.estado)}</span></td><td>${dinero(t.saldo_posterior)}</td></tr>`).join(''):'<tr><td colspan="5" class="empty">Aún no hay transacciones registradas.</td></tr>'}catch(e){$('mensaje').hidden=false;$('mensaje').className='message bad';$('mensaje').textContent=e.message}}
function elegir(tipo){operacion=tipo;$('abono').classList.toggle('selected',tipo==='ABONO');$('retiro').classList.toggle('selected',tipo==='RETIRO');$('enviar').firstChild.textContent=tipo==='ABONO'?'Realizar abono ':'Realizar retiro '}
$('abono').addEventListener('click',()=>elegir('ABONO'));
$('retiro').addEventListener('click',()=>elegir('RETIRO'));
$('refrescar').addEventListener('click',actualizar);
$('form').addEventListener('submit',async ev=>{ev.preventDefault();const boton=$('enviar');boton.disabled=true;try{const t=await obtener('/api/transacciones',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({operacion,monto:Number($('monto').value)})});$('mensaje').hidden=false;$('mensaje').className='message '+(t.estado==='EXITO'?'good':'bad');$('mensaje').textContent=`${t.estado}: ${t.detalle}. Saldo: ${dinero(t.saldo_posterior)}`;$('monto').value='';await actualizar()}catch(e){$('mensaje').hidden=false;$('mensaje').className='message bad';$('mensaje').textContent=e.message}finally{boton.disabled=false}});
actualizar();