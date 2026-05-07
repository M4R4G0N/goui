package components

import (
	"fmt"
	"strings"
	"time"
)

type CalendarRangeComponent struct {
	ID    string
	Name  string
	Class string
	Value string // O valor atual (ex: "01/01/2024 - 15/01/2024")
}

func (c *CalendarRangeComponent) String() string {
	return c.Value
}

func (c *CalendarRangeComponent) GetID() string {
	if c.ID == "" {
		c.ID = AutoID()
	}
	Register(c.ID, c) // Se identifica para o framework
	return c.ID
}

// CalendarRange creates a button that opens a range selection calendar.
func CalendarRange(opts ...any) *CalendarRangeComponent {
	c := &CalendarRangeComponent{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			c.ID = v
		case Name:
			c.Name = string(v)
		case ID:
			c.ID = string(v)
		}
	}
	return c
}

func (c *CalendarRangeComponent) Render() string {
	id := c.GetID()
	jsID := strings.ReplaceAll(id, "-", "_")
	now := time.Now()
	
	months := []string{"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho", "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"}
	monthDrop := Dropdown(ID("month-"+id), months, Value(months[now.Month()-1])).Render()
	
	years := []string{}
	for y := now.Year() - 20; y <= now.Year()+20; y++ {
		years = append(years, fmt.Sprintf("%d", y))
	}
	yearDrop := Dropdown(ID("year-"+id), years, Value(fmt.Sprintf("%d", now.Year()))).Render()

	return fmt.Sprintf(`
		<div class="goui-calendar-wrapper" style="position: relative;">
			<button type="button" class="goui-input goui-calendar-btn" id="btn-%s" onclick="window.gouiToggleCalendar_%s()">
				Selecionar datas
			</button>
			
			<div id="pop-%s" class="goui-calendar-popup" style="display: none; position: absolute; top: 100%%; left: 0; z-index: 1000; margin-top: 5px;">
				<div class="goui-calendar-header">
					<button type="button" onclick="window.gouiPrevMonth_%s()">&#9664;</button>
					<div style="width: 120px;">%s</div>
					<div style="width: 90px;">%s</div>
					<button type="button" onclick="window.gouiNextMonth_%s()">&#9654;</button>
				</div>
				<div class="goui-calendar-grid weekdays">
					<div class="goui-calendar-day-head">Dom</div>
					<div class="goui-calendar-day-head">Seg</div>
					<div class="goui-calendar-day-head">Ter</div>
					<div class="goui-calendar-day-head">Qua</div>
					<div class="goui-calendar-day-head">Qui</div>
					<div class="goui-calendar-day-head">Sex</div>
					<div class="goui-calendar-day-head">Sáb</div>
				</div>
				<div class="goui-calendar-grid" id="days-%s"></div>
				<div class="goui-calendar-footer" id="footer-%s">Selecionado: -</div>
			</div>
			<!-- Dois inputs com o mesmo nome para enviar como slice no Go -->
			<input type="hidden" name="%s" id="%s-start">
			<input type="hidden" name="%s" id="%s-end">
		</div>
		<script>
			(function() {
				const id = '%s';
				const jsID = '%s';
				let startDate = null;
				let endDate = null;
				let currDate = new Date();
				const months = ['Janeiro','Fevereiro','Março','Abril','Maio','Junho','Julho','Agosto','Setembro','Outubro','Novembro','Dezembro'];

				window['gouiToggleCalendar_' + jsID] = function() {
					const pop = document.getElementById('pop-' + id);
					if(pop.style.display === 'none') {
						const sVal = document.getElementById(id + '-start').value;
						const eVal = document.getElementById(id + '-end').value;
						if(sVal && eVal) {
							const p1 = sVal.split('/');
							const p2 = eVal.split('/');
							startDate = new Date(p1[2], p1[1]-1, p1[0]);
							endDate = new Date(p2[2], p2[1]-1, p2[0]);
							currDate = new Date(startDate);
							updateHeader();
						}
						window['gouiRenderCalendar_' + jsID]();
						pop.style.display = 'block';
					} else {
						pop.style.display = 'none';
					}
				};
				
				function updateHeader() {
					const mInp = document.getElementById('month-' + id);
					const yInp = document.getElementById('year-' + id);
					if(mInp) gouiSelectOption('month-' + id, months[currDate.getMonth()], months[currDate.getMonth()]);
					if(yInp) gouiSelectOption('year-' + id, currDate.getFullYear().toString(), currDate.getFullYear().toString());
				}

				window['gouiPrevMonth_' + jsID] = () => { currDate.setMonth(currDate.getMonth()-1); updateHeader(); window['gouiRenderCalendar_' + jsID](); };
				window['gouiNextMonth_' + jsID] = () => { currDate.setMonth(currDate.getMonth()+1); updateHeader(); window['gouiRenderCalendar_' + jsID](); };

				window['gouiRenderCalendar_' + jsID] = function() {
					const monthInp = document.getElementById('month-' + id);
					const yearInp = document.getElementById('year-' + id);
					const daysCont = document.getElementById('days-' + id);
					if(!monthInp || !yearInp || !daysCont) return;

					if(!monthInp.dataset.listened) {
						monthInp.addEventListener('change', () => { 
							const mIdx = months.indexOf(monthInp.value);
							currDate.setMonth(mIdx !== -1 ? mIdx : parseInt(monthInp.value)); 
							window['gouiRenderCalendar_' + jsID](); 
						});
						yearInp.addEventListener('change', () => { 
							currDate.setFullYear(parseInt(yearInp.value)); 
							window['gouiRenderCalendar_' + jsID](); 
						});
						monthInp.dataset.listened = 'true';
					}

					const y = parseInt(yearInp.value || currDate.getFullYear());
					const m = parseInt(monthInp.value ? (months.indexOf(monthInp.value) !== -1 ? months.indexOf(monthInp.value) : monthInp.value) : currDate.getMonth());
					const firstDay = new Date(y, m, 1).getDay();
					const totalDays = new Date(y, m + 1, 0).getDate();

					daysCont.innerHTML = '';
					for(let i=0; i<firstDay; i++) daysCont.innerHTML += '<div class="goui-calendar-day-empty"></div>';
					
					for(let d=1; d<=totalDays; d++) {
						const date = new Date(y, m, d);
						const div = document.createElement('div');
						div.className = 'goui-calendar-day';
						div.innerText = d;
						
						const time = date.getTime();
						if(startDate && time === startDate.getTime()) div.classList.add('selected');
						if(endDate && time === endDate.getTime()) div.classList.add('selected');
						if(startDate && endDate && date > startDate && date < endDate) div.classList.add('in-range');
						if(new Date().toDateString() === date.toDateString()) div.classList.add('today');

						div.onclick = (e) => {
							e.stopPropagation();
							if(!startDate || (startDate && endDate)) {
								startDate = date; endDate = null;
							} else if(date < startDate) {
								endDate = startDate; startDate = date;
								document.getElementById('pop-' + id).style.display = 'none';
							} else {
								endDate = date;
								document.getElementById('pop-' + id).style.display = 'none';
							}
							updateDisplay(id);
							window['gouiRenderCalendar_' + jsID]();
						};
						daysCont.appendChild(div);
					}
				};

				function updateDisplay(id) {
					const btn = document.getElementById('btn-' + id);
					const foot = document.getElementById('footer-' + id);
					const hStart = document.getElementById(id + '-start');
					const hEnd = document.getElementById(id + '-end');
					
					const fmt = (d) => d.toLocaleDateString('pt-BR');
					
					if(startDate && endDate) {
						const txt = fmt(startDate) + ' - ' + fmt(endDate);
						if(btn) btn.innerText = txt;
						if(foot) foot.innerText = 'Selecionado: ' + txt;
						if(hStart) hStart.value = fmt(startDate);
						if(hEnd) hEnd.value = fmt(endDate);
						// Dispatch event on first input to trigger reactivity
						if(hStart) hStart.dispatchEvent(new Event('change', { bubbles: true }));
					} else if(startDate) {
						if(btn) btn.innerText = fmt(startDate) + '...';
						if(foot) foot.innerText = 'Selecionado: ' + fmt(startDate);
						if(hStart) hStart.value = fmt(startDate);
						if(hEnd) hEnd.value = '';
					}
				}

				window.addEventListener('click', (e) => {
					if (!e.target.closest('.goui-calendar-wrapper')) {
						const pop = document.getElementById('pop-' + id);
						if(pop) pop.style.display = 'none';
					}
				});
			})();
		</script>
	`, id, jsID, id, jsID, monthDrop, yearDrop, jsID, id, id, c.Name, id, c.Name, id, id, jsID)
}

// DateRange is an alias for CalendarRange to maintain backward compatibility with existing code.
func DateRange(opts ...any) *CalendarRangeComponent {
	return CalendarRange(opts...)
}
