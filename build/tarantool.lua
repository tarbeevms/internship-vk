box.cfg{
    -- параметры конфигурации, такие как listen, memtx_memory, etc.
    listen = 3301
}


box.schema.user.create('root', {password = 'root'})
box.schema.user.grant('root', 'read,write,execute', 'universe')

-- Функция для создания спейсов
local function create_spaces()

    -- Спейс для сессий
    if not box.space.sessions then
        box.schema.create_space('sessions', {
            format = {
                {name = 'username', type = 'string'},
                {name = 'token', type = 'string'}
            }
        })

        -- Уникальный индекс на поле username
        box.space.sessions:create_index('primary', {
            type = 'hash',
            parts = {'username'}
        })

        -- Уникальный индекс на поле token (token_index)
        box.space.sessions:create_index('token_index', {
        type = 'hash',
        parts = {'token'}
        })
    end

    -- Спейс для пользователей
    if not box.space.users then
        box.schema.create_space('users', {
            format = {
                {name = 'username', type = 'string'},
                {name = 'password', type = 'string'}
            }
        })

        -- Уникальный индекс на поле username
        box.space.users:create_index('primary', {
            type = 'hash',
            parts = {'username'}
        })

        -- Вставляем запись с пользователем admin и паролем presale
        -- Проверяем, существует ли уже запись для admin, чтобы избежать дублирования
        if not box.space.users:get('admin') then
            box.space.users:insert{'admin', 'presale'}
            print("Admin user added.")
        else
            print("Admin user already exists.")
        end
        
    end

    -- Спейс для данных (ключ-значение)
    if not box.space.data then
        box.schema.create_space('data', {
            format = {
                {name = 'key', type = 'string'},
                {name = 'value', type = 'scalar'}
            }
        })

        -- Уникальный индекс на поле key
        box.space.data:create_index('primary', {
            type = 'hash',
            parts = {'key'}
        })
    end
end

-- Инициализация спейсов
create_spaces()

print("Spaces 'sessions', 'users' and 'data' initialized.")
