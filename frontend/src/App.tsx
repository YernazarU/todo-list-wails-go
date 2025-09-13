import {useEffect, useMemo, useState} from 'react';
import './App.css';
import {CreateTask, DeleteTask, ListTasks, UpdateTask} from "../wailsjs/go/main/App";
import {repo} from "../wailsjs/go/models";

// Типы для приоритетов задач
type Priority = 'low'|'medium'|'high'

// Алиас для типа задачи из Go бэкенда
type Task = repo.Task

/**
 * Главный компонент приложения TODO
 * Управляет состоянием задач и взаимодействием с бэкендом
 */
function App() {
  // Состояние списка задач
  const [tasks, setTasks] = useState<Task[]>([])
  
  // Состояние формы добавления новой задачи
  const [title, setTitle] = useState('')
  const [dueAt, setDueAt] = useState('')
  const [priority, setPriority] = useState<Priority>('medium')
  
  // Состояние для модального окна подтверждения удаления
  const [confirmId, setConfirmId] = useState<string|undefined>()
  
  // Состояние для отображения ошибок пользователю
  const [error, setError] = useState<string|undefined>()

  /**
   * Обновляет список задач с сервера
   * Вызывается при загрузке компонента и после изменений
   */
  async function refresh() {
    try {
      console.log('Refreshing tasks...')
      const list = await ListTasks()
      console.log('Tasks received:', list)
      setTasks(list)
    } catch (error) {
      console.error('Error refreshing tasks:', error)
      setError('Ошибка загрузки задач: ' + (error as Error).message)
    }
  }

  // Загружаем задачи при монтировании компонента
  useEffect(() => { 
    refresh()
    
    // Добавляем тестовые задачи для демонстрации функциональности
    // Используем setTimeout чтобы не блокировать первоначальную загрузку
    setTimeout(async () => {
      try {
        await CreateTask('Изучить Wails', '', '', 'high')
        await CreateTask('Создать TODO приложение', '', '', 'medium')
        await CreateTask('Получить работу', '', '', 'high')
        await refresh()
      } catch (error) {
        console.log('Test tasks already exist or error:', error)
      }
    }, 1000)
  }, [])

  /**
   * Добавляет новую задачу
   * Валидирует ввод и отправляет данные на сервер
   */
  async function addTask() {
    try {
      const trimmed = title.trim()
      if (!trimmed) return // Не добавляем пустые задачи
      
      console.log('Adding task:', { title: trimmed, dueAt, priority })
      
      // Конвертируем дату в ISO формат для бэкенда
      const dueISO = dueAt ? new Date(dueAt).toISOString() : ''
      
      await CreateTask(trimmed, '', dueISO, priority)
      console.log('Task created successfully')
      
      // Очищаем форму после успешного добавления
      setTitle(''); setDueAt(''); setPriority('medium')
      await refresh()
    } catch (error) {
      console.error('Error adding task:', error)
      setError('Ошибка добавления задачи: ' + (error as Error).message)
    }
  }

  /**
   * Переключает статус задачи (выполнена/не выполнена)
   * @param task - задача для переключения
   */
  async function toggle(task: Task) {
    try {
      console.log('Toggling task:', task.id, 'completed:', !task.completed)
      
      // Создаем копию задачи с измененным статусом
      const updated = new repo.Task({...task, completed: !task.completed})
      
      await UpdateTask(updated)
      await refresh()
    } catch (error) {
      console.error('Error toggling task:', error)
    }
  }

  /**
   * Удаляет задачу по ID
   * @param id - идентификатор задачи для удаления
   */
  async function remove(id: string) {
    try {
      console.log('Removing task:', id)
      await DeleteTask(id)
      
      // Закрываем модальное окно подтверждения
      setConfirmId(undefined)
      await refresh()
    } catch (error) {
      console.error('Error removing task:', error)
    }
  }

  // Разделяем задачи на активные и выполненные для отображения
  const active = tasks.filter(t => !t.completed)
  const completed = tasks.filter(t => t.completed)

  return (
    <div id="app" style={{maxWidth: 820, margin: '0 auto', padding: 16}}>
      <h1>To‑Do List</h1>
      <div style={{fontSize: 12, opacity: 0.7, marginBottom: 16}}>
        Всего задач: {tasks.length}, Активных: {active.length}, Выполненных: {completed.length}
      </div>
      
      {error && (
        <div style={{background: '#ff6b6b', color: 'white', padding: 8, borderRadius: 4, marginBottom: 16}}>
          {error}
          <button onClick={() => setError(undefined)} style={{marginLeft: 8, background: 'none', border: 'none', color: 'white', cursor: 'pointer'}}>×</button>
        </div>
      )}

      <div style={{display:'grid', gap:8, gridTemplateColumns:'1fr 160px 140px'}}>
        <input placeholder="Новая задача" value={title} onChange={e=>setTitle(e.target.value)} />
        <input type="datetime-local" value={dueAt} onChange={e=>setDueAt(e.target.value)} />
        <select value={priority} onChange={e=>setPriority(e.target.value as Priority)}>
          <option value="low">Низкий</option>
          <option value="medium">Средний</option>
          <option value="high">Высокий</option>
        </select>
      </div>
      <div style={{marginTop:8, display: 'flex', gap: 8}}>
        <button className="btn" onClick={addTask} disabled={!title.trim()}>Добавить</button>
        <button className="btn" onClick={async () => {
          try {
            for (const task of tasks) {
              await DeleteTask(task.id!)
            }
            await refresh()
          } catch (error) {
            console.error('Error clearing tasks:', error)
          }
        }} style={{background: '#ff6b6b'}}>Очистить все</button>
      </div>

      <section style={{marginTop:16}}>
        <h2>Активные</h2>
        <ul style={{listStyle:'none', padding:0, margin:0, display:'grid', gap:8}}>
          {active.map(t => (
            <li key={t.id} style={{display:'grid', gridTemplateColumns:'24px 1fr auto', gap:8, alignItems:'center', padding:8, border:'1px solid #ddd', borderRadius:8}}>
              <input type="checkbox" checked={false} onChange={()=>toggle(t)} />
              <div>
                <div style={{fontWeight:600}}>{t.title}</div>
                <div style={{fontSize:12, opacity:0.8}}>
                  {t.priority === 1 ? 'НИЗКИЙ' : t.priority === 2 ? 'СРЕДНИЙ' : 'ВЫСОКИЙ'} {t.due_date?`• до ${new Date(t.due_date as string).toLocaleString()}`:''}
                </div>
              </div>
              <div style={{display:'flex', gap:8}}>
                <button className="btn" onClick={()=>setConfirmId(t.id!)}>Удалить</button>
              </div>
            </li>
          ))}
        </ul>
      </section>

      <section style={{marginTop:24}}>
        <h2>Выполненные</h2>
        <ul style={{listStyle:'none', padding:0, margin:0, display:'grid', gap:8}}>
          {completed.map(t => (
            <li key={t.id} style={{display:'grid', gridTemplateColumns:'24px 1fr auto', gap:8, alignItems:'center', padding:8, border:'1px solid #ddd', borderRadius:8, opacity:0.7}}>
              <input type="checkbox" checked={true} onChange={()=>toggle(t)} />
              <div>
                <div style={{textDecoration:'line-through'}}>{t.title}</div>
                <div style={{fontSize:12, opacity:0.8}}>
                  {t.priority === 1 ? 'НИЗКИЙ' : t.priority === 2 ? 'СРЕДНИЙ' : 'ВЫСОКИЙ'} {t.due_date?`• до ${new Date(t.due_date as string).toLocaleString()}`:''}
                </div>
              </div>
              <div style={{display:'flex', gap:8}}>
                <button className="btn" onClick={()=>setConfirmId(t.id!)}>Удалить</button>
              </div>
            </li>
          ))}
        </ul>
      </section>

      {confirmId && (
        <div style={{position:'fixed', inset:0, background:'rgba(0,0,0,0.3)', display:'grid', placeItems:'center'}}>
          <div style={{background:'#fff', color:'#000', padding:16, borderRadius:8, minWidth:280}}>
            <div style={{marginBottom:12}}>Удалить задачу?</div>
            <div style={{display:'flex', gap:8, justifyContent:'flex-end'}}>
              <button className="btn" onClick={()=>setConfirmId(undefined)}>Отмена</button>
              <button className="btn" onClick={()=>remove(confirmId!)}>Удалить</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default App
