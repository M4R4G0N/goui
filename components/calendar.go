package components

import (
	"fmt"
	"strings"
	"time"
)

type CalendarComponent struct {
	ID    string
	Name  string
	Value string
	Class string
	Style Style
}

func (c *CalendarComponent) String() string {
	return c.Value
}

func (c *CalendarComponent) GetID() string {
	if c.ID == "" {
		c.ID = AutoID()
	}
	Register(c.ID, c) // Se identifica para o framework
	return c.ID
}

// Calendar creates a custom-styled date picker.
func Calendar(opts ...any) *CalendarComponent {
	c := &CalendarComponent{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			ParseStringAttr(v, &c.Class, &c.ID, nil)
		case Name:
			c.Name = string(v)
		case ID:
			c.ID = string(v)
		case Value:
			c.Value = string(v)
		}
	}
	return c
}

func (c *CalendarComponent) Render() string {
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
		<div class="goui-calendar-wrapper" style="position: relative;" id="wrapper-%s">
			<div class="goui-input-icon-wrapper">
				<input type="text" id="%s" name="%s" class="goui-input goui-calendar-input" 
					placeholder="DD/MM/AAAA" value="%s" readonly onclick="window.gouiToggleCalendar_%s()">
				<span class="goui-input-icon">
					<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
				</span>
			</div>
			
			<div id="pop-%s" class="goui-calendar-popup" style="display: none; position: absolute; top: 100%%; left: 0; z-index: 1000; margin-top: 5px;">
				<div class="goui-calendar-header">
					<button type="button" onclick="window.gouiPrevMonth_%s()">&#9664;</button>
					<div style="width: 120px;">%s</div>
					<div style="width: 90px;">%s</div>
					<button type="button" onclick="window.gouiNextMonth_%s()">&#9654;</button>
				</div>
				<div class="goui-calendar-grid weekdays">
					<div class="goui-calendar-day-head">D</div>
					<div class="goui-calendar-day-head">S</div>
					<div class="goui-calendar-day-head">T</div>
					<div class="goui-calendar-day-head">Q</div>
					<div class="goui-calendar-day-head">Q</div>
					<div class="goui-calendar-day-head">S</div>
					<div class="goui-calendar-day-head">S</div>
				</div>
				<div class="goui-calendar-grid" id="days-%s"></div>
			</div>
		</div>
		<script>
			(function() {
				const id = '%s';
				const jsID = '%s';
				let currDate = new Date();
				const months = ['Janeiro','Fevereiro','Março','Abril','Maio','Junho','Julho','Agosto','Setembro','Outubro','Novembro','Dezembro'];

				window['gouiToggleCalendar_' + jsID] = function() {
					const pop = document.getElementById('pop-' + id);
					if(pop.style.display === 'none') {
						const val = document.getElementById(id).value;
						if(val) {
							const parts = val.split('/');
							if(parts.length === 3) {
								currDate = new Date(parts[2], parts[1]-1, parts[0]);
								updateHeader();
							}
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
					const mainInp = document.getElementById(id);
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
					
					const selVal = mainInp.value;
					for(let d=1; d<=totalDays; d++) {
						const date = new Date(y, m, d);
						const div = document.createElement('div');
						div.className = 'goui-calendar-day';
						div.innerText = d;
						
						const fmtDate = date.toLocaleDateString('pt-BR');
						if(selVal === fmtDate) div.classList.add('selected');
						if(new Date().toDateString() === date.toDateString()) div.classList.add('today');

						div.onclick = (e) => {
							e.stopPropagation();
							mainInp.value = fmtDate;
							document.getElementById('pop-' + id).style.display = 'none';
							mainInp.dispatchEvent(new Event('change', { bubbles: true }));
						};
						daysCont.appendChild(div);
					}
				};

				window.addEventListener('click', (e) => {
					if (!e.target.closest('.goui-calendar-wrapper')) {
						const pop = document.getElementById('pop-' + id);
						if(pop) pop.style.display = 'none';
					}
				});
			})();
		</script>
	`, id, id, id, c.Value, jsID, id, jsID, monthDrop, yearDrop, jsID, id, id, jsID)
}
