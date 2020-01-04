import styled from 'styled-components';
import { space, border, flexbox, typography } from 'styled-system';

const Group = styled.div`
    display: flex;
    flex-direction: column;

    &:last-child {
        margin-bottom: 0;
    }

    ${space} ${border} ${flexbox} ${typography}
`;

Group.defaultProps = {
  marginBottom: 5,
};

export default Group;
