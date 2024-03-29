CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$
                                     
    DECLARE                         
        data json;
        notification json;
                                                    
    BEGIN                                        
                           
        -- Convert the old or new row to JSON, based on the kind of action.                                        
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE                           
            data = row_to_json(NEW);
        END IF;
        
        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
        
                        
        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('traccar.events',notification::text);
        
        -- Result is ignored since this is an AFTER trigger
        RETURN NULL; 
    END;
    
$$ LANGUAGE plpgsql;


CREATE TRIGGER tc_positions_notify_event
AFTER INSERT ON public.tc_positions
    FOR EACH ROW EXECUTE PROCEDURE notify_event();

CREATE TRIGGER tc_events_notify_event
AFTER INSERT ON public.tc_events
    FOR EACH ROW EXECUTE PROCEDURE notify_event();

CREATE TRIGGER tc_devices_notify_event
AFTER INSERT OR UPDATE OR DELETE ON public.tc_devices
    FOR EACH ROW EXECUTE PROCEDURE notify_event();